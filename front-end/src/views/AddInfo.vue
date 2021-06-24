<template>
  <div> 
    <!-- Date picker --> 
     

    <!-- Add files --> 
    <Add name="Productos"/>
    <Add name="Compradores"/>
    <Add name="Transacciones"/>

    <!-- Progress icon --> 
    <v-progress-circular
    indeterminate
    size="80"
    width="10"
    color="primary"
    id="progress"
    ></v-progress-circular>

    <v-alert
      id="alert"
      type="error"
    > Los recuadros no pueden estar vacíos al momento de agregar la información </v-alert>

  </div>
</template>

<script>

  import Add from '../components/Add.vue'

  export default {
    name: 'Home',

    components: {
      Add,
    }
  }

</script>

<style scoped>
  #progress{
    display: none; 
    position: absolute;
    top: 50%;
    left: 50%;  
  }
  #alert {
    display: none; 
    position: absolute;
    bottom: 50%; 
  }
</style>