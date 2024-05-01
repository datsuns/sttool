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
  function onTextConfigChanged(event, type) {
    //LogPrint(`Notify fired! Detail: ${event.detail.value}`);
    switch (type) {
      case "obsip":
        Config.ObsIp = event.detail.value;
        break;
      case "obspass":
        Config.ObsPass = event.detail.value;
        break;
      default:
        LogPrint(`onTextConfigChanged: invalid type: ${type}`);
        return;
    }
    issueDispatch(Config);
  }
  function onNumberConfigChanged(event, type) {
    LogPrint(`Notify fired! Detail: ${event.detail.value}`);
    let v = Number(event.detail.value);
    if (isNaN(v)) {
      LogPrint(`invalid number text: ${event.detail.value}`);
      return;
    }
    switch (type) {
      case "width":
        Config.ClipPlayerWidth = v;
        break;
      case "height":
        Config.ClipPlayerHeight = v;
        break;
      case "port":
        Config.LocalServerPortNumber = v;
        break;
      case "obsport":
        Config.ObsPort = v;
        break;
      default:
        LogPrint(`onNumberConfigChanged invalid type: ${type}`);
        return;
    }
    issueDispatch(Config);
  }
</script>

<h1>設定画面</h1>
<Paper>
  <Content>オーバーレイ設定</Content>
  <BoolConfig
    value={Config.OverlayEnabled}
    labelText="オーバーレイ有効"
    on:changed={onOverlayConfigChanged}
  ></BoolConfig>
  <Paper square variant="outlined">
    <Content>URL</Content>
    <Title>http://localhost:{Config.LocalServerPortNumber}</Title>
    <TextConfig
      value={Config.LocalServerPortNumber}
      labelText="port番号"
      valueType="number"
      on:changed={(e) => onNumberConfigChanged(e, "port")}
    ></TextConfig>
  </Paper>
  <Paper square variant="outlined">
    <Title>クリップ再生サイズ</Title>
    <TextConfig
      value={Config.ClipPlayerWidth}
      labelText="幅"
      valueType="number"
      on:changed={(e) => onNumberConfigChanged(e, "width")}
    ></TextConfig>
    <TextConfig
      value={Config.ClipPlayerHeight}
      labelText="高さ"
      valueType="number"
      on:changed={(e) => onNumberConfigChanged(e, "height")}
    ></TextConfig>
  </Paper>
</Paper>

<Paper>
  <Content>OBS連携</Content>
  <Paper square variant="outlined">
    <Content>websocketサーバー</Content>
    <TextConfig
      value={Config.ObsIp}
      labelText="サーバーIP"
      valueType="text"
      on:changed={(e) => onTextConfigChanged(e, "obsip")}
    ></TextConfig>
    <TextConfig
      value={Config.ObsPort}
      labelText="サーバーポート"
      valueType="number"
      on:changed={(e) => onNumberConfigChanged(e, "obsport")}
    ></TextConfig>
    <TextConfig
      value={Config.ObsPass}
      labelText="サーバーパスワード"
      valueType="text"
      on:changed={(e) => onTextConfigChanged(e, "obspass")}
    ></TextConfig>
  </Paper>
</Paper>

<Paper>
  <StaticTextConfig
    value={Config.NotifySoundFile}
    labelText="新規クリップ通知音"
    selectionFilter="audio/mp3, audio/wav"
    on:changed={onClipNotificationChanged}
  ></StaticTextConfig>
</Paper>
