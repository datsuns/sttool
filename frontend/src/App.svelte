<script>
  import { onMount } from "svelte";
  import { writable } from "svelte/store";
  import Drawer from "@smui/drawer";
  import TopAppBar from "@smui/top-app-bar";
  import IconButton, { Icon } from "@smui/icon-button";
  import List, { Item } from "@smui/list";
  import { LoadConfig, SaveConfig } from "../wailsjs/go/main/App.js";
  import { LogPrint, EventsOn } from "../wailsjs/runtime/runtime";
  import MainScreen from "./MainScreen.svelte";
  import ConfigScreen from "./ConfigScreen.svelte";

  let drawerOpened = false;
  let currentScreen = writable("main");

  let mainScreenRef;

  let Clips = [];
  let Config;
  let Debug = false;

  onMount(() => {
    LoadConfig().then((result) => {
      Config = result;
      Debug = Config.DebugMode;
      LogPrint(`App:onMount debug : ${Debug}`);
    });
  });

  function toggleDrawer() {
    drawerOpened = !drawerOpened;
  }

  function switchScreen(screen) {
    LogPrint(`App:switchScreen >> ${screen}`);
    currentScreen.set(screen);
    toggleDrawer();
    LogPrint(`App:switchScreen << ${screen}`);
  }

  EventsOn("OnConnected", (msg) => {
    LogPrint(`App:OnConnected ${msg}`);
    mainScreenRef.handleOnConnected(msg);
  });

  EventsOn("OnRaid", (msg, username, items) => {
    LogPrint(`App:OnRaid ${msg}`);
    let entry = { name: username, body: items };
    LogPrint(`raid from ${username}`);
    items.forEach((c) => {
      LogPrint(`user clip [${c.Title}], url [${c.Thumbnail}], mp4 [${c.Mp4}]`);
    });
    Clips = [...Clips, entry];
  });

  function onConfigChanged(event) {
    let cfg = event.detail.value;
    //LogPrint(`cfg changed ${cfg.OverlayEnabled}`);
    SaveConfig(cfg);
  }
</script>

<main>
  <!-- SMUI Styles -->
  <link
    rel="stylesheet"
    href="/src/smui.css"
    media="(prefers-color-scheme: light)"
  />
  <link
    href="https://fonts.googleapis.com/css2?family=Material+Icons&display=swap"
    rel="stylesheet"
  />

  <TopAppBar variant="fixed" collapsed>
    <IconButton on:click={toggleDrawer} class="material-icons">
      <Icon class="material-icons">menu</Icon>
    </IconButton>
  </TopAppBar>

  <Drawer class="app-drawer" variant="dismissible" bind:open={drawerOpened}>
    <List>
      <Item on:click={() => switchScreen("main")}>
        <span class="smui-list-item__text">メイン</span>
      </Item>
      <Item on:click={() => switchScreen("settings")}>
        <span class="smui-list-item__text">設定</span>
      </Item>
    </List>
  </Drawer>

  <div class="content">
    {#if $currentScreen === "main"}
      <MainScreen
        bind:this={mainScreenRef}
        raidUserClips={Clips}
        debugMode={Debug}
      />
    {:else if $currentScreen === "settings"}
      <ConfigScreen {Config} on:changed={onConfigChanged} />
    {/if}
  </div>
</main>

<style>
  :global(.mdc-card) {
    background-color: rgba(18, 29, 45, 1);
  }
</style>
