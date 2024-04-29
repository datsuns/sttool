<script>
    import { createEventDispatcher } from "svelte";
    import FormField from "@smui/form-field";
    import Checkbox from "@smui/checkbox";
    import Paper from "@smui/paper";
    import Textfield from "@smui/textfield";
    import { LogPrint } from "../wailsjs/runtime/runtime";

    const dispatch = createEventDispatcher();

    export let value;
    export let labelText;
    export let valueType;

    function handleChanged(event) {
        LogPrint(`TextConfig ${labelText} changed to ${event.target.value}`);
        dispatch("changed", {
            value: event.target.value,
        });
    }
</script>

<link rel="stylesheet" href="/src/style.css" />
<Textfield
    bind:value
    on:change={handleChanged}
    label={labelText}
    type={valueType}
></Textfield>
