<template>
    <form @submit.prevent="onSubmit">
        <h2 class="label-wrapper">
        <label for="new-todo-input">What needs to be done </label>
        <input
            type="text"
            id="new-todo-input"
            name="new-todo"
            autocomplete="off"
            v-model.lazy.trim="label"/>
        <button type="submit" class="btn btn__primary btn__lg">Add</button>
        </h2>
    </form>
</template>

<script>

    export default{
        methods:{
            onSubmit(){
                if(this.label===""){
                    return;
                }
                this.$emit("todo-added",this.label);  //event name is todo-added
                this.label=""
            }
        },
        data(){
            return{
                label:"",
            }
        }
    };
</script>