<script>
  import { onMount } from "svelte";
  import { GetServerPort } from "../wailsjs/go/main/App.js";
  import { LogPrint, EventsOn } from "../wailsjs/runtime/runtime";
  import MainScreen from "./MainScreen.svelte";

  let mainScreenRef;

  let Clips = [];
  let ServerPort = 0;

  onMount(() => {
    GetServerPort().then((result) => (ServerPort = result));
  });

  EventsOn("OnConnected", (msg) => {
    LogPrint(`App:OnConnected ${msg}`);
    mainScreenRef.handleOnConnected(msg);
  });

  EventsOn("OnRaid", (msg, username, items) => {
    LogPrint(`App:OnRaid ${msg}`);
    mainScreenRef.handleOnCRaid(msg, username, items);
  });
</script>

<main>
  <!-- SMUI Styles -->
  <link
    rel="stylesheet"
    href="/src/smui.css"
    media="(prefers-color-scheme: light)"
  />
  <MainScreen
    bind:this={mainScreenRef}
    overlayServerPort={ServerPort}
    raidUserClips={Clips}
  />
</main>

<style>
  :global(.mdc-card) {
    background-color: rgba(18, 29, 45, 1);
  }
</style>
