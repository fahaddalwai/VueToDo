<template>
    <div class="stack-small" v-if="!isEditing">
      <div class="custom-checkbox">
        <input
          type="checkbox"
          class="checkbox"
          :id="id"
          :checked="isDone"
          @change="$emit('checkbox-changed')" />
        <label :for="id" class="checkbox-label">{{label}}</label>
      </div>
      <div class="btn-group">
        <button type="button" class="btn" @click="toggleToItemEditForm">
          Edit <span class="visually-hidden">{{label}}</span>
        </button>
        <button type="button" class="btn btn__danger" @click="deleteToDo">
          Delete <span class="visually-hidden">{{label}}</span>
        </button>
      </div>
      
    </div>
    
    <to-do-item-edit-form v-else :id="id" :label="label"></to-do-item-edit-form>

  </template>

<script>
import ToDoItemEditForm from "./ToDoItemEditForm";
export default{
    components: {
  ToDoItemEditForm
},
    props:{ //props are how you pass down data from parent to component(1 way)
        label:{required:true,type:String},
        done:{default:false,type:Boolean},
        id:{required:true,type:String}
    },
    data(){ //Used to store and alter variables(2 way)
        return{
            isDone:this.done, //isDone variable takes the value of the done prop
            isEditing: false
        }
    },
    methods: {
    deleteToDo() {
      this.$emit('item-deleted');
    },
    toggleToItemEditForm() {
      this.isEditing = true;
    },
    itemEdited(newLabel){
      this.$emit('item-edited', newLabel);
      this.isEditing = false;
    },
    editCancelled() {
    this.isEditing = false;
    }
  }
};
</script>