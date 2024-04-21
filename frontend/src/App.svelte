<script>
  import { onMount } from "svelte";
  import LayoutGrid, { Cell } from "@smui/layout-grid";
  import Card, {
    Content,
    PrimaryAction,
    Media,
    MediaContent,
  } from "@smui/card";
  import {
    OpenURL,
    StartClip,
    StopClip,
    GetServerPort,
  } from "../wailsjs/go/main/App.js";
  import { DebugRaidTest } from "../wailsjs/go/main/App.js";
  import { LogPrint, EventsOn } from "../wailsjs/runtime/runtime";
  import Clip from "./Clip.svelte";

  let dbg_RaidUser = "";
  let Clips = [];
  let ServerPort = 0;

  onMount(() => {
    GetServerPort().then((result) => (ServerPort = result));
  });

  function openLink(url) {
    OpenURL(url).then((result) => LogPrint("open link"));
  }

  const callDebugRaidTest = async () => {
    await DebugRaidTest(dbg_RaidUser);
  };

  function startClipTest(url, duration) {
    StartClip(url, duration).then((result) => LogPrint("Clip finished"));
  }

  const stopClipTest = async () => {
    await StopClip();
    LogPrint("stop Clip");
  };

  EventsOn("OnConnected", (msg) => {
    LogPrint(`OnConnected ${msg}`);
    DebugRaidTest("datsuns7");
  });

  EventsOn("OnRaid", (msg, username, items) => {
    let entry = { name: username, body: items };
    LogPrint(`raid from ${username}`);
    items.forEach((c) => {
      LogPrint(`user clip [${c.Title}], url [${c.Thumbnail}], mp4 [${c.Mp4}]`);
    });
    Clips = [...Clips, entry];
  });
</script>

<main>
  <!-- SMUI Styles -->
  <link
    rel="stylesheet"
    href="/src/smui.css"
    media="(prefers-color-scheme: light)"
  />
  <input
    bind:value={dbg_RaidUser}
    class="input"
    placeholder="レイドテスト用(ユーザID)"
  />
  <button on:click={callDebugRaidTest}>raid test</button>
  <div class="my-overlay-url">
    http://localhost:{ServerPort}
  </div>
  <button on:click={stopClipTest}>stop clip</button>
  {#each Clips.slice().reverse() as clip}
    <h1>{clip.name} さんのクリップ</h1>
    <LayoutGrid>
      {#each clip.body as c}
        <Cell span={4}>
          <div style="height: 100%;">
            <Clip
              startClipCallback={startClipTest}
              Url={c.Mp4}
              Title={c.Title}
              Thumnail={c.Thumbnail}
              Duration={c.Duration}
              ViewCount={c.ViewCount}
            />
          </div>
        </Cell>
      {/each}
    </LayoutGrid>
  {/each}
</main>

<style>
  :global(.mdc-card) {
    background-color: rgba(18, 29, 45, 1);
  }

  .my-overlay-url {
    color: lightblue;
  }
</style>
