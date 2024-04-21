<script>
    import { LogPrint } from "../wailsjs/runtime/runtime";
    import Card, {
        Content,
        PrimaryAction,
        Media,
        MediaContent,
    } from "@smui/card";

    export let startClipCallback;
    export let Title;
    export let Thumnail;
    export let Url;
    export let Duration;
    export let ViewCount;

    const issueClipPlay = async () => {
        await startClipCallback();
        LogPrint("stop Clip");
    };
</script>

<Card>
    <PrimaryAction
        style="height: 100%;"
        on:click={() => startClipCallback(Url, Duration)}
    >
        <Media class="card-media-16x9" aspectRatio="16x9">
            <MediaContent>
                <img class="my-thumbnail" src={Thumnail} alt={Title} />
            </MediaContent>
        </Media>
        <Content>
            <div class="my-clip-title">{Title}</div>
        </Content>
        <Content>
            <div class="my-viewcount">
                再生数:{ViewCount} / 時間:{Duration}
            </div>
        </Content>
    </PrimaryAction>
</Card>

<style>
    :global(.mdc-card) {
        background-color: rgba(18, 29, 45, 1);
    }

    .my-thumbnail {
        width: 100%;
        height: 100%;
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
