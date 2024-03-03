<script>
  import Button, { Label } from '@smui/button';
  import { OpenURL } from '../wailsjs/go/main/App.js';
  import { DebugRaidTest  } from '../wailsjs/go/main/App.js';
  import { LogPrint } from '../wailsjs/runtime/runtime'

  let clicked = 0;
  let dbg_RaidUser = "";

  function openLink(url) {
    OpenURL(url).then(result => LogPrint("Opened2"));
  };

  const callDebugRaidTest = async () => {
    await DebugRaidTest(dbg_RaidUser);
  };

</script>

<main>
  <link rel="stylesheet" href="node_modules/svelte-material-ui/bare.css" />
  <pre class="status">Clicked: {clicked}</pre>
  <Button on:click={() => clicked++} variant="raised">
    <Label>Raised</Label>
  </Button>
  <Button color="secondary" on:click={() => clicked++} variant="raised">
    <Label>Raised</Label>
  </Button>
  <input bind:value={dbg_RaidUser}  class="input" placeholder="debug raid user" />
  <button class="btn" on:click={() => openLink('https://www.google.com/')}>Greet</button>
  <button on:click={callDebugRaidTest}>raid test</button>
</main>

<style>
  .btn {
    width: 60px;
    height: 30px;
    line-height: 30px;
    border-radius: 3px;
    border: none;
    margin: 0 0 0 20px;
    padding: 0 8px;
    cursor: pointer;
  }

  .btn:hover {
    background-image: linear-gradient(to top, #cfd9df 0%, #e2ebf0 100%);
    color: #333333;
  }

</style>
