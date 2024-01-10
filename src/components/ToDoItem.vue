<template>
    <div>
        <input type="checkbox" :id="id" :checked="isDone"/>
        <label :for="id">{{label}}</label>
    </div>
</template>

<script>

export default{
    props:{ //props are how you pass down data from parent to component(1 way)
        label:{required:true,type:String},
        done:{default:false,type:Boolean},
        id:{required:true,type:String}
    },
    data(){ //Used to store and alter variables(2 way)
        return{
            isDone:this.done, //isDone variable takes the value of the done prop
        }
    }
};
</script>