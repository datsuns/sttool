<script>
    import { createEventDispatcher } from "svelte";
    import FormField from "@smui/form-field";
    import Checkbox from "@smui/checkbox";
    import Paper from "@smui/paper";
    import { LogPrint } from "../wailsjs/runtime/runtime";

    const dispatch = createEventDispatcher();

    export let value;
    export let labelText;

    function handleChanged(event) {
        //LogPrint(`BoolConfig ${labelText} changed to ${event.target.checked}`);
        dispatch("changed", {
            checked: event.target.checked,
        });
    }
</script>

<link rel="stylesheet" href="/src/style.css" />
<Paper square variant="outlined">
    <FormField>
        <Checkbox bind:checked={value} on:change={handleChanged} />
        <span class="my-custom-tilte" slot="label">{labelText}</span>
    </FormField>
</Paper>
