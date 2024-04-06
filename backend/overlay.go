package backend

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type OverlayContext struct {
	ChannStartClip   chan struct{}
	ChannStopClip    chan struct{}
	ServerPort       int
	PlayMarginSecond int
	ClipId           string
}

const OverlayHtml = `
<!DOCTYPE html>
<html>
<head>
    <title>Go Server-Sent Events Example</title>
</head>
<body>
    <div id="iframe-container"></div>

    <script>
        const evtSource = new EventSource("/events");

        evtSource.addEventListener("on", function(event) {
            const data = JSON.parse(event.data);
            const container = document.getElementById('iframe-container');
            container.innerHTML = ` +
	"`<iframe src=\"${data.src}\" width=\"640\" height=\"480\"></iframe>`;" +
	`
        });

        evtSource.addEventListener("off", function(event) {
            const container = document.getElementById('iframe-container');
            container.innerHTML = ''; // iframeをクリア
        });
    </script>
</body>
</html>
`

func NewOverlay(cfg *Config) *OverlayContext {
	ret := &OverlayContext{
		ChannStartClip:   make(chan struct{}),
		ChannStopClip:    make(chan struct{}),
		ServerPort:       cfg.LocalServerPortNumber,
		PlayMarginSecond: cfg.ClipPlayIntervalMarginSecond,
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

func (o *OverlayContext) buildPlayInterval(clipDuration float32) time.Duration {
	margin := time.Second * time.Duration(o.PlayMarginSecond)
	interval := time.Second*time.Duration(clipDuration) + margin
	return interval
}

func (o *OverlayContext) StartClip(clipId string, duration float32) {
	o.ClipId = clipId
	o.ChannStartClip <- struct{}{}
	interval := o.buildPlayInterval(duration)
	logger.Info("Overlay:interval", slog.Any("time", interval))
	time.Sleep(interval)
	if o.ClipId != clipId {
		logger.Info("Overlay", slog.Any("msg", "other clip already started"))
		return
	}
	o.ClipId = ""
	o.ChannStopClip <- struct{}{}
	logger.Info("Overlay:finished")
}

func (o *OverlayContext) StopClip() {
	o.ClipId = ""
	o.ChannStopClip <- struct{}{}
	logger.Info("Overlay:forceStop")
}

func (o *OverlayContext) OnEvent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	select {
	case <-o.ChannStartClip:
		src := buildSrcUrl(o.ClipId)
		fmt.Fprintf(w, "event: on\ndata: {\"src\": \"%s\"}\n\n", src)
		logger.Info("Ovelay:ON", slog.Any("clip ID", o.ClipId), slog.Any("URL", src))
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
	go func() { o.Main(cfg.LocalServerPortNumber) }()
}
