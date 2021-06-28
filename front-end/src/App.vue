<!---------------------------------------------->
<!--              Template                    --> 
<!---------------------------------------------->
<template>
  <v-app>

    <!---------------------------------------------->
    <!--               Navbar                     --> 
    <!----------------------------------------------> 
    <v-app-bar app dark>
      <!-- Title --> 
      <v-toolbar-title class="ml-5"> <span class="T">T</span> <i class="ienda">ienda</i> </v-toolbar-title>
      <v-spacer></v-spacer>

      <!-- Navigation links --> 
      <router-link to="/" class="mr-10 navigation"> Agregar </router-link>
      <router-link to="/buyers" class="mr-5 navigation"> Compradores </router-link>

    </v-app-bar>

    <!---------------------------------------------->
    <!--               Router views               --> 
    <!---------------------------------------------->
    <v-main class="mt-5">
      <router-view/>
    </v-main>

    <!---------------------------------------------->
    <!--                 Footer                   -->
    <!---------------------------------------------->
    <v-footer padless>
        <v-col class="text-center" cols="12" >
          Made with ♡ by JD
      </v-col>
    </v-footer>

  </v-app>
</template>

<!---------------------------------------------->
<!--               Script                     --> 
<!---------------------------------------------->
<script>

  export default {
    name: 'App',
  };

</script>

<!---------------------------------------------->
<!--                Style                     --> 
<!---------------------------------------------->
<style>
  /*----------------- Background color -----------------*/
  #app {
    background-color: #190862 !important;
  }
  .v-app-bar {
    background-color: #190862 !important;
  }

  /*----------------- Navbar -----------------*/
  /* Title */
  .T {
    font-family: 'Lobster', cursive;
    font-size: 40px;
  }
  .ienda {
    font-size: 40px;
    margin-left: -10px !important;
  }
  /* Navigation */
  .navigation {
    text-decoration: none;
    color: #ffffff !important; 
    font-size: 1.2rem;
  }

</style>