<script>
  import { LogPrint } from "../wailsjs/runtime/runtime";
  // ここで必要なスクリプトを追加
  export let Config;

  const callDebugRaidTest = async () => {
    LogPrint("load");
  };
</script>

<div>
  <button on:click={callDebugRaidTest}>raid test</button>
  <h1>設定画面</h1>
  <p>ここに設定関連のコンテンツが入ります {Config.NotifySoundFile}</p>
  <p>ここに設定関連のコンテンツが入ります {Config.OverlayEnabled}</p>
</div>
