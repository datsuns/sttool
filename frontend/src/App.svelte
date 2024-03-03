<script>
  import logo from './assets/images/logo-universal.png'
  import {Greet} from '../wailsjs/go/main/App.js'
  import Button, { Label } from '@smui/button';
  import { OpenURL } from '../wailsjs/go/main/App.js';
  import { DebugRaidTest  } from '../wailsjs/go/main/App.js';
  import { LogPrint } from '../wailsjs/runtime/runtime'

  let resultText = "Please enter your name below 👇"
  let name
  let clicked = 0;
  let dbg_RaidUser = "";

  function openLink(url) {
    OpenURL(url).then(result => LogPrint("Opened2"));
  };

  const callDebugRaidTest = async () => {
    await DebugRaidTest(dbg_RaidUser);
  };

  function greet() {
    LogPrint("Greet");
    Greet(name).then(result => resultText = result)
  }
</script>

<main>
  <link rel="stylesheet" href="node_modules/svelte-material-ui/bare.css" />
  <img alt="Wails logo" id="logo" src="{logo}">
  <div class="result" id="result">{resultText}</div>
  <pre class="status">Clicked: {clicked}</pre>
  <Button on:click={() => clicked++} variant="raised">
    <Label>Raised</Label>
  </Button>
  <Button color="secondary" on:click={() => clicked++} variant="raised">
    <Label>Raised</Label>
  </Button>
  <div class="input-box" id="input">
    <input autocomplete="off" bind:value={name} class="input" id="name" type="text"/>
    <button class="btn" on:click={() => openLink('https://www.google.com/')}>Greet</button>
  </div>
  <input bind:value={dbg_RaidUser}  class="input" placeholder="debug raid user" />
  <button on:click={callDebugRaidTest}>raid test</button>
</main>

<style>

  #logo {
    display: block;
    width: 50%;
    height: 50%;
    margin: auto;
    padding: 10% 0 0;
    background-position: center;
    background-repeat: no-repeat;
    background-size: 100% 100%;
    background-origin: content-box;
  }

  .result {
    height: 20px;
    line-height: 20px;
    margin: 1.5rem auto;
  }

  .input-box .btn {
    width: 60px;
    height: 30px;
    line-height: 30px;
    border-radius: 3px;
    border: none;
    margin: 0 0 0 20px;
    padding: 0 8px;
    cursor: pointer;
  }

  .input-box .btn:hover {
    background-image: linear-gradient(to top, #cfd9df 0%, #e2ebf0 100%);
    color: #333333;
  }

  .input-box .input {
    border: none;
    border-radius: 3px;
    outline: none;
    height: 30px;
    line-height: 30px;
    padding: 0 10px;
    background-color: rgba(240, 240, 240, 1);
    -webkit-font-smoothing: antialiased;
  }

  .input-box .input:hover {
    border: none;
    background-color: rgba(255, 255, 255, 1);
  }

  .input-box .input:focus {
    border: none;
    background-color: rgba(255, 255, 255, 1);
  }

</style>
