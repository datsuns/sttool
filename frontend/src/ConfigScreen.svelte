<script>
  import { createEventDispatcher } from "svelte";
  import { LogPrint } from "../wailsjs/runtime/runtime";
  import BoolConfig from "./BoolConfig.svelte";
  import StaticTextConfig from "./StaticTextConfig.svelte";
  import { onMount } from "svelte";
  export let Config;

  const dispatch = createEventDispatcher();

  onMount(() => {
    //LogPrint(`config: ${Config.NotifySoundFile}`);
    //LogPrint(`config: ${Config.OverlayEnabled}`);
  });

  function onOverlayConfigChanged(event) {
    Config.OverlayEnabled = event.detail.checked;
    dispatch("changed", {
      value: Config,
    });
  }
  function callbackFunction2(event) {
    LogPrint(`Notify fired! Detail: ${event.detail.value}`);
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
  on:changed={callbackFunction2}
></StaticTextConfig>
