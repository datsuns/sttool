<script>
  import Button, { Label } from "@smui/button";
  import LayoutGrid, { Cell } from "@smui/layout-grid";
  import Card, {
    Content,
    PrimaryAction,
    Media,
    MediaContent,
    Actions,
    ActionButtons,
    ActionIcons,
  } from "@smui/card";
  import { OpenURL } from "../wailsjs/go/main/App.js";
  import { DebugRaidTest } from "../wailsjs/go/main/App.js";
  import { LogPrint, EventsOn } from "../wailsjs/runtime/runtime";

  let dbg_RaidUser = "";
  let Clips = [];

  function openLink(url) {
    OpenURL(url).then((result) => LogPrint("Opened2"));
  }

  const callDebugRaidTest = async () => {
    await DebugRaidTest(dbg_RaidUser);
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
  <link rel="stylesheet" href="node_modules/svelte-material-ui/bare.css" />
  <input
    bind:value={dbg_RaidUser}
    class="input"
    placeholder="レイドテスト用(ユーザID)"
  />
  <button on:click={callDebugRaidTest}>raid test</button>
  {#each Clips.slice().reverse() as clip}
    <h1>{clip.name} さんのクリップ</h1>
    <LayoutGrid>
      {#each clip.body as c}
        <Cell span="3" desktop="3">
          <Card>
            <PrimaryAction on:click={() => openLink(c.Url)}>
              <Media class="card-media-16x9" aspectRatio="16x9">
                <MediaContent>
                  <img class="Thumbnail" src={c.Thumbnail} alt={c.Title} />
                </MediaContent>
              </Media>
              <Content class="mdc-typography--body2" style="color: blue;">
                <p>{c.Title}</p>
              </Content>
              <Content class="mdc-typography--body2" style="color: black;">
                再生数:{c.ViewCount}
              </Content>
            </PrimaryAction>
          </Card>
        </Cell>
      {/each}
    </LayoutGrid>
  {/each}
</main>

<style>
  .Thumbnail {
    width: 75%;
    height: 100%;
  }
</style>
