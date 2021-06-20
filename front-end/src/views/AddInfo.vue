<template>
  <div> 
    <!-- Date picker --> 
    

    <!-- Add files --> 
    <Add title="Productos"/>
    <Add title="Compradores"/>
    <Add title="Transacciones"/>
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