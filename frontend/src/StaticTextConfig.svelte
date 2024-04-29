<script>
    import { createEventDispatcher } from "svelte";
    import Button, { Label } from "@smui/button";
    import { LogPrint } from "../wailsjs/runtime/runtime";
    import Paper, { Title, Content } from "@smui/paper";

    const dispatch = createEventDispatcher();

    export let value = "";
    export let labelText = "";
    export let selectionFilter = "";
    let fileInput;

    function handleChanged(event) {
        //LogPrint(`on change ${event.target.value}`);
        value = event.target.value;
        dispatch("changed", {
            value: event.target.value,
        });
    }
    function handleClick() {
        fileInput.click();
    }
</script>

<link rel="stylesheet" href="/src/style.css" />

<input
    accept={selectionFilter}
    bind:this={fileInput}
    id="avatar"
    name="avatar"
    type="file"
    on:input={handleChanged}
/>
<Paper square variant="outlined">
    <Content>{labelText}</Content>
    <Title>{value}</Title>
    <Button color="secondary" on:click={handleClick} variant="raised">
        <Label>変更</Label>
    </Button>
</Paper>

<style>
    input[type="file"] {
        display: none; /* 元のinput要素を非表示にする */
    }
</style>
