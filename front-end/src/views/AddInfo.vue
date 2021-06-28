<!---------------------------------------------->
<!--              Template                    --> 
<!---------------------------------------------->
<template>
  <div> 
    
    <!------------------------------------->
    <!--            Date picker          --> 
    <!------------------------------------->
    <div class="date-picker-wrapper mb-5 mt-2 text-center">
      <label class="title white--text">
        Selecciona una fecha : 
        <input type="date" id="selected-date"/> 
      </label>
    </div>

    <!------------------------------------->
    <!--            Add files            --> 
    <!------------------------------------->
    <Add name="Productos"/>
    <Add name="Compradores"/>
    <Add name="Transacciones"/>

    <!------------------------------------->
    <!--          Progress icon          --> 
    <!------------------------------------->
    <v-progress-circular
    indeterminate
    size="80"
    width="10"
    color="primary"
    id="progress"
    ></v-progress-circular>

    <!------------------------------------->
    <!--            Error alert          -->
    <!------------------------------------->
    <v-alert
      id="alert"
      type="error"
    > Los recuadros no pueden estar vacíos al momento de agregar la información </v-alert>

  </div>
</template>

<!---------------------------------------------->
<!--               Script                     --> 
<!---------------------------------------------->
<script>

  import Add from '../components/Add.vue'

  export default {
    name: 'Home',

    components: {
      Add,
    }
  }

</script>

<!---------------------------------------------->
<!--                Style                     --> 
<!---------------------------------------------->
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
  .date-picker-wrapper {
    width: 100%;  
  }
  .date-picker-wrapper input {
    background-color: #F3F6FF !important;
    border-radius: 10px
  }
</style>