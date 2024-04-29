<script>
  import { createEventDispatcher } from "svelte";
  import Paper, { Title, Content } from "@smui/paper";
  import { LogPrint } from "../wailsjs/runtime/runtime";
  import BoolConfig from "./BoolConfig.svelte";
  import StaticTextConfig from "./StaticTextConfig.svelte";
  import TextConfig from "./TextConfig.svelte";
  import { onMount } from "svelte";
  export let Config;

  const dispatch = createEventDispatcher();

  onMount(() => {
    //LogPrint(`config: ${Config.NotifySoundFile}`);
    //LogPrint(`config: ${Config.OverlayEnabled}`);
  });

  function issueDispatch(cfg) {
    dispatch("changed", {
      value: cfg,
    });
  }

  function onOverlayConfigChanged(event) {
    Config.OverlayEnabled = event.detail.checked;
    issueDispatch(Config);
  }
  function onClipNotificationChanged(event) {
    Config.NotifySoundFile = event.detail.value;
    //LogPrint(`Notify fired! Detail: ${event.detail.value}`);
    issueDispatch(Config);
  }
  function onClipWidthChanged(event) {
    LogPrint(`Notify fired! Detail: ${event.detail.value}`);
    let v = Number(event.detail.value);
    if (isNaN(v)) {
      LogPrint(`invalid number text: ${event.detail.value}`);
      return;
    }
    Config.ClipPlayerWidth = v;
    issueDispatch(Config);
  }
  function onClipHeightChanged(event) {
    LogPrint(`Notify fired! Detail: ${event.detail.value}`);
    let v = Number(event.detail.value);
    if (isNaN(v)) {
      LogPrint(`invalid number text: ${event.detail.value}`);
      return;
    }
    Config.ClipPlayerHeight = v;
    issueDispatch(Config);
  }
</script>

<h1>設定画面</h1>
<BoolConfig
  value={Config.OverlayEnabled}
  labelText="オーバーレイ有効"
  on:changed={onOverlayConfigChanged}
></BoolConfig>
<StaticTextConfig
  value={Config.NotifySoundFile}
  labelText="新規クリップ通知音"
  selectionFilter="audio/mp3, audio/wav"
  on:changed={onClipNotificationChanged}
></StaticTextConfig>

<Paper square variant="outlined">
  <Title>オーバーレイクリップサイズ</Title>
  <TextConfig
    value={Config.ClipPlayerWidth}
    labelText="幅"
    valueType="number"
    on:changed={onClipWidthChanged}
  ></TextConfig>
  <TextConfig
    value={Config.ClipPlayerHeight}
    labelText="高さ"
    valueType="number"
    on:changed={onClipHeightChanged}
  ></TextConfig>
</Paper>
