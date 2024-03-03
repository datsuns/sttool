<script>
  import Button, { Label } from '@smui/button';
  import LayoutGrid, { Cell } from '@smui/layout-grid';
  import Card, { Content, PrimaryAction, Media, MediaContent, Actions, ActionButtons, ActionIcons, } from '@smui/card';
  import { OpenURL } from '../wailsjs/go/main/App.js';
  import { DebugRaidTest  } from '../wailsjs/go/main/App.js';
  import { LogPrint, EventsOn } from '../wailsjs/runtime/runtime'

  let clicked = 0;
  let dbg_RaidUser = "";
  let Clips = [];

  function openLink(url) {
    OpenURL(url).then(result => LogPrint("Opened2"));
  };

  const callDebugRaidTest = async () => {
    await DebugRaidTest(dbg_RaidUser);
  };

  EventsOn("OnRaid", (msg, username, items) => {
    let entry = { name: username, body: items}
    LogPrint(`raid from ${username}`)
    items.forEach(c => {
        LogPrint(`user clip [${c.Title}]`);
    });
    Clips = [...Clips, entry];
  });

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
  {#each Clips as clip}
    <h1>{clip.name} さんのクリップ</h1>
    <LayoutGrid>
      {#each clip.body as c}
        <Card>
          <PrimaryAction on:click={() => openLink(c.Url)}>
            <Media aspectRatio="16x9">
              <img src={c.Thumbnail} alt={c.Title}/>
            </Media>
            <Content class="mdc-typography--body2" style="color: blue;">
              {c.Title}
            </Content>
          </PrimaryAction>
        </Card>
      {/each}
    </LayoutGrid>
  {/each}
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
