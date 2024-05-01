<script>
    import { writable } from "svelte/store";
    import { onMount } from "svelte";
    import LayoutGrid, { Cell } from "@smui/layout-grid";
    import Paper, { Title } from "@smui/paper";
    import {
        OpenURL,
        StartClip,
        StopClip,
        StopObsStream,
        DebugRaidTest,
    } from "../wailsjs/go/main/App.js";
    import { LogPrint, EventsOn } from "../wailsjs/runtime/runtime";
    import Clip from "./Clip.svelte";

    export let raidUserClips = [];
    export let debugMode = false;

    let Debug = writable(false);
    let dbg_RaidUser = "";

    onMount(() => {
        if (debugMode) {
            Debug.set(true);
        }
    });

    function openLink(url) {
        OpenURL(url).then((result) => LogPrint("open link"));
    }

    const callDebugRaidTest = async () => {
        await DebugRaidTest(dbg_RaidUser);
    };

    function startClip(url, duration) {
        StartClip(url, duration).then((result) => LogPrint("Clip finished"));
    }

    const stopClip = async () => {
        await StopClip();
        LogPrint("stop Clip");
    };

    const stopStream = async () => {
        await StopObsStream();
        LogPrint("stream stopped");
    };

    export function handleOnConnected(msg) {
        LogPrint(`MainScreen:handleOnConnected ${msg}`);
        if (debugMode) {
            DebugRaidTest("datsuns7");
        }
    }
</script>

{#if $Debug}
    <input
        bind:value={dbg_RaidUser}
        class="input"
        placeholder="レイドテスト用(ユーザID)"
    />
    <button on:click={callDebugRaidTest}>raid test</button>
    <button on:click={stopStream}>配信停止</button>
{/if}
<button on:click={stopClip}>クリップ強制停止</button>
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
                            startClipCallback={startClip}
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
</style>
