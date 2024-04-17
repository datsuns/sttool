package backend

import (
	"fmt"
	"log/slog"
	"net/http"
)

type OverlayContext struct {
	ChannStartClip   chan struct{}
	ChannStopClip    chan struct{}
	ServerPort       int
	PlayMarginSecond int
	ClipUrl          string
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
	"`<video id=\"clip-player-body\" autoplay width=\"640\" height=\"480\"> <source src=\"${data.src}\"> </video>`;" +
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
		ServerPort:     cfg.LocalPortNum(),
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

func (o *OverlayContext) OnEvent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	select {
	case <-o.ChannStartClip:
		//src := buildSrcUrl(o.ClipId)
		src := o.ClipUrl
		fmt.Fprintf(w, "event: on\ndata: {\"src\": \"%s\"}\n\n", src)
		logger.Info("Ovelay:ON", slog.Any("clip ID", o.ClipUrl), slog.Any("URL", src))
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

func (o *OverlayContext) Main(serverPort int) {
	logger.Info("Ovelay:Start")
	http.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		o.OnEvent(w, r)
	})
	http.HandleFunc("/", rootDocument)

	if err := http.ListenAndServe(
		fmt.Sprintf(":%v", serverPort),
		nil,
	); err != nil {
		logger.Error("Ovelay:ERROR", slog.Any("error", err.Error))

	}
	logger.Info("Ovelay:Finish")
}

func (o *OverlayContext) Serve(cfg *Config) {
	go func() { o.Main(cfg.LocalPortNum()) }()
}
