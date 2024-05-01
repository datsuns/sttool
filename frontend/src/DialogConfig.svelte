<script>
    import { createEventDispatcher } from "svelte";
    import Button, { Label } from "@smui/button";
    import { LogPrint } from "../wailsjs/runtime/runtime";
    import Paper, { Title, Content } from "@smui/paper";
    import {
        OpenFileDialog,
        OpenDiectoryDialog,
    } from "../wailsjs/go/main/App.js";

    const dispatch = createEventDispatcher();

    export let type = "";
    export let value = "";
    export let labelText = "";
    export let selectionFilter = "";

    function handleChanged(v) {
        LogPrint(`on change ${v}`);
        dispatch("changed", {
            value: v,
        });
    }
    function handleClick(t) {
        LogPrint(`handleClick ${t}`);
        switch (t) {
            case "file":
                LogPrint(`>> OpenFileDialog ${value}, ${selectionFilter}`);
                OpenFileDialog(value, selectionFilter).then((result) => {
                    if (result === "") {
                        return;
                    }
                    handleChanged(result);
                });
                break;
            case "dir":
                OpenDiectoryDialog(value).then((result) => {
                    if (result === "") {
                        return;
                    }
                    handleChanged(result);
                });
                break;
            default:
                LogPrint(`handleClick: invalid type: ${t}`);
                return;
        }
    }
</script>

<link rel="stylesheet" href="/src/style.css" />

<Paper square variant="outlined">
    <Content>{labelText}</Content>
    <Title>{value}</Title>
    <Button color="secondary" on:click={handleClick(type)} variant="raised">
        <Label>変更</Label>
    </Button>
</Paper>

<style>
</style>
