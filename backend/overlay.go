package backend

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type OverlayContext struct {
	ChannStartClip   chan struct{}
	ChannStopClip    chan struct{}
	PlayMarginSecond int
	ClipUrl          string
	ServeMux         *http.ServeMux
	Server           *http.Server
}

const OverlayHtml = `
<!DOCTYPE html>
<html>
<head>
    <title>Go Server-Sent Events Example</title>
</head>
<body>
    <div id="clip-player"></div>

    <script>
        const evtSource = new EventSource("/events");
        const container = document.getElementById('clip-player');

        evtSource.addEventListener("on", function(event) {
            console.log('start play');
            const data = JSON.parse(event.data);
            container.innerHTML = ` +
	"`<video id=\"clip-player-body\" autoplay width=\"${data.width}\" height=\"${data.height}\"> <source src=\"${data.src}\"> </video>`;" +
	`
            const player = document.getElementById('clip-player-body');
            player.addEventListener('ended', function () {
                console.log('play ended');
                container.innerHTML = ''; // iframeをクリア
            });
            player.play();
        });

        evtSource.addEventListener("off", function(event) {
            console.log('force stop');
            container.innerHTML = ''; // iframeをクリア
        });

    </script>
</body>
</html>
`

func NewOverlay(cfg *Config) *OverlayContext {
	ret := &OverlayContext{
		ChannStartClip: make(chan struct{}),
		ChannStopClip:  make(chan struct{}),
	}
	return ret
}

func rootDocument(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, OverlayHtml)
	logger.Info("Ovelay:rootDocument")
}

func buildSrcUrl(clipID string) string {
	return fmt.Sprintf(
		"https://clips.twitch.tv/embed?clip=%v&parent=localhost&autoplay=true&muted=false",
		clipID,
	)
}

func (o *OverlayContext) StartClip(url string, duration float32) {
	o.ClipUrl = url
	o.ChannStartClip <- struct{}{}
	logger.Info("Overlay:StartClip", slog.Any("url", url), slog.Any("clip time", duration))
}

func (o *OverlayContext) StopClip() {
	if o.ClipUrl == "" {
		logger.Info("Overlay", slog.Any("msg", "clip already stopped"))
		return
	}
	o.ChannStopClip <- struct{}{}
	logger.Info("Overlay:forceStop")
}

func (o *OverlayContext) OnEvent(w http.ResponseWriter, r *http.Request, cfg *Config) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	select {
	case <-o.ChannStartClip:
		//src := buildSrcUrl(o.ClipId)
		src := o.ClipUrl
		fmt.Fprintf(w,
			"event: on\ndata: {\"src\": \"%s\", \"width\": \"%v\", \"height\": \"%v\"}\n\n",
			src, cfg.ClipWidth(), cfg.ClipHeight(),
		)
		logger.Info("Ovelay:ON",
			slog.Any("clip ID", o.ClipUrl),
			slog.Any("URL", src),
			slog.Any("x", cfg.ClipWidth()),
			slog.Any("y", cfg.ClipHeight()),
		)
	case <-o.ChannStopClip:
		fmt.Fprintf(w, "event: off\ndata: {}\n\n")
		logger.Info("Ovelay:OFF")
	}
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	} else {
		logger.Error("Ovelay:Streaming unsupported.")
	}

}

func (o *OverlayContext) Main(cfg *Config) {
	logger.Info("Ovelay:Start")
	o.ServeMux = http.NewServeMux()
	o.ServeMux.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		o.OnEvent(w, r, cfg)
	})
	o.ServeMux.HandleFunc("/", rootDocument)
	o.Server = &http.Server{
		Addr:    fmt.Sprintf(":%v", cfg.LocalPortNum()),
		Handler: o.ServeMux,
	}

	if err := o.Server.ListenAndServe(); err != nil {
		logger.Error("Ovelay:ERROR", slog.Any("error", err.Error))
	}
	logger.Info("Ovelay:Finish")
}

func (o *OverlayContext) Serve(cfg *Config) {
	go func() {
		o.Main(cfg)
	}()
}

func (o *OverlayContext) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	o.Server.Shutdown(ctx)
	<-ctx.Done()
}
