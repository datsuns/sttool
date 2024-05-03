<script>
  import { createEventDispatcher } from "svelte";
  import Paper, { Title, Content } from "@smui/paper";
  import Button, { Label } from "@smui/button";
  import Snackbar, { Actions } from "@smui/snackbar";
  import IconButton from "@smui/icon-button";
  import { LogPrint } from "../wailsjs/runtime/runtime";
  import BoolConfig from "./BoolConfig.svelte";
  import DialogConfig from "./DialogConfig.svelte";
  import TextConfig from "./TextConfig.svelte";
  import { onMount } from "svelte";
  import { TestObsConnection } from "../wailsjs/go/main/App.js";

  export let Config;
  let showObsConnectionResult;
  let ObsConnectionResultBody = "";

  const dispatch = createEventDispatcher();

  onMount(() => {
    //LogPrint(`config: ${Config.NotifySoundFile}`);
    //LogPrint(`config: ${Config.OverlayEnabled}`);
  });

  function testObsConnection() {
    TestObsConnection().then((result) => {
      LogPrint(`testObsConnection [${result}]`);
      if (result.length > 0) {
        ObsConnectionResultBody = `接続OK OBSバージョン[${result}]`;
      } else {
        ObsConnectionResultBody = "接続エラー";
      }
      showObsConnectionResult.open();
    });
  }

  function issueDispatch(cfg) {
    dispatch("changed", {
      value: cfg,
    });
  }

  function onBoolConfigChanged(event, type) {
    switch (type) {
      case "overlayen":
        Config.OverlayEnabled = event.detail.checked;
        break;
      default:
        LogPrint(`onBoolConfigChanged: invalid type: ${type}`);
        return;
    }
    issueDispatch(Config);
  }

  function onTextConfigChanged(event, type) {
    //LogPrint(`Notify fired! Detail: ${event.detail.value}`);
    switch (type) {
      case "logdest":
        Config.LogDest = event.detail.value;
        break;
      case "obsip":
        Config.ObsIp = event.detail.value;
        break;
      case "obspass":
        Config.ObsPass = event.detail.value;
        break;
      case "clipsound":
        Config.NotifySoundFile = event.detail.value;
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
  <DialogConfig
    type="dir"
    value={Config.LogDest}
    labelText="ログ出力フォルダ"
    on:changed={(e) => onTextConfigChanged(e, "logdest")}
  ></DialogConfig>
</Paper>

<Paper>
  <Title>オーバーレイ設定</Title>
  <BoolConfig
    value={Config.OverlayEnabled}
    labelText="オーバーレイ有効"
    on:changed={(e) => onBoolConfigChanged(e, "overlayen")}
  ></BoolConfig>
  <Paper square variant="outlined">
    <Content>URL</Content>
    <Content>http://localhost:{Config.LocalServerPortNumber}</Content>
    <TextConfig
      value={Config.LocalServerPortNumber}
      labelText="port番号"
      valueType="number"
      on:changed={(e) => onNumberConfigChanged(e, "port")}
    ></TextConfig>
  </Paper>
  <Paper square variant="outlined">
    <Content>クリップ再生サイズ</Content>
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
      valueType="password"
      on:changed={(e) => onTextConfigChanged(e, "obspass")}
    ></TextConfig>
    <Button color="secondary" on:click={testObsConnection} variant="raised">
      <Label>接続テスト</Label>
    </Button>
    <Snackbar bind:this={showObsConnectionResult}>
      <Label>{ObsConnectionResultBody}</Label>
      <Actions>
        <IconButton class="material-icons" title="Dismiss">close</IconButton>
      </Actions>
    </Snackbar>
  </Paper>
</Paper>

<Paper>
  <DialogConfig
    type="file"
    value={Config.NotifySoundFile}
    labelText="新規クリップ通知音"
    selectionFilter="*.mp3; *.wav"
    on:changed={(e) => onTextConfigChanged(e, "clipsound")}
  ></DialogConfig>
</Paper>
