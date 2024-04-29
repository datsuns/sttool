<script>
    import LayoutGrid, { Cell } from "@smui/layout-grid";
    import Paper, { Title } from "@smui/paper";
    import {
        OpenURL,
        StartClip,
        StopClip,
        DebugRaidTest,
    } from "../wailsjs/go/main/App.js";

    import { LogPrint, EventsOn } from "../wailsjs/runtime/runtime";

    import Clip from "./Clip.svelte";

    export let overlayServerPort;
    export let raidUserClips = [];

    let dbg_RaidUser = "";

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

    export function handleOnConnected(msg, debugMode) {
        LogPrint(`MainScreen:handleOnConnected ${msg}`);
        if (debugMode) {
            DebugRaidTest("datsuns7");
        }
    }
</script>

<input
    bind:value={dbg_RaidUser}
    class="input"
    placeholder="レイドテスト用(ユーザID)"
/>

<button on:click={callDebugRaidTest}>raid test</button>
<div class="my-overlay-url">
    http://localhost:{overlayServerPort}
</div>
<button on:click={stopClipTest}>stop clip</button>
{#each raidUserClips.slice().reverse() as clip}
    <h1>{clip.name} さんのクリップ</h1>
    {#if clip.body.length == 0}
        <Paper>
            <Title>クリップがありません</Title>
        </Paper>
    {:else}
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
    {/if}
{/each}

<style>
    :global(.mdc-card) {
        background-color: rgba(18, 29, 45, 1);
    }

    .my-overlay-url {
        color: lightblue;
    }
</style>
