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

  let dbg_RaidUser = "";
  let Clips = [];
  let ServerPort = 0;

  onMount(() => {
    //DebugRaidTest("twitch-id");
    GetServerPort().then((result) => (ServerPort = result));
  });

  function openLink(url) {
    OpenURL(url).then((result) => LogPrint("open link"));
  }

  const callDebugRaidTest = async () => {
    await DebugRaidTest(dbg_RaidUser);
  };

  function startClipTest(id, duration) {
    StartClip(id, duration).then((result) => LogPrint("Clip finished"));
  }

  const stopClipTest = async () => {
    await StopClip();
    LogPrint("stop Clip");
  };

  EventsOn("OnRaid", (msg, username, items) => {
    let entry = { name: username, body: items };
    LogPrint(`raid from ${username}`);
    items.forEach((c) => {
      LogPrint(`user clip [${c.Title}], url [${c.Thumbnail}]`);
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
          <Card style="height: 100%;">
            <PrimaryAction
              style="height: 100%;"
              on:click={() => startClipTest(c.Id, c.Duration)}
            >
              <Media class="card-media-16x9" aspectRatio="16x9">
                <MediaContent>
                  <img class="my-thumbnail" src={c.Thumbnail} alt={c.Title} />
                </MediaContent>
              </Media>
              <Content>
                <div class="my-clip-title">{c.Title}</div>
              </Content>
              <Content>
                <div class="my-viewcount">再生数:{c.ViewCount}</div>
              </Content>
            </PrimaryAction>
          </Card>
        </Cell>
      {/each}
    </LayoutGrid>
  {/each}
</main>

<style>
  :global(.mdc-card) {
    background-color: rgba(18, 29, 45, 1);
  }

  .my-thumbnail {
    width: 100%;
    height: 100%;
  }

  .my-overlay-url {
    color: lightblue;
  }

  .my-clip-title {
    color: whitesmoke;
  }

  .my-viewcount {
    position: absolute;
    color: darkgray;
    bottom: 8px;
    left: 0px;
    right: 0px;
    font-size: 70%;
  }
</style>
