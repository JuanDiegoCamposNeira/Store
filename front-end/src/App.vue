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

    data: () => ({
      //
    }),
  };

</script>

<!---------------------------------------------->
<!--                Style                     --> 
<!---------------------------------------------->
<style>
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

<!-- 
<div class="d-flex align-center">
        <v-img
          alt="Vuetify Logo"
          class="shrink mr-2"
          contain
          src="https://cdn.vuetifyjs.com/images/logos/vuetify-logo-dark.png"
          transition="scale-transition"
          width="40"
        />

        <v-img
          alt="Vuetify Name"
          class="shrink mt-1 hidden-sm-and-down"
          contain
          min-width="100"
          src="https://cdn.vuetifyjs.com/images/logos/vuetify-name-dark.png"
          width="100"
        />
      </div>

      <v-spacer></v-spacer>

      <v-btn
        href="https://github.com/vuetifyjs/vuetify/releases/latest"
        target="_blank"
        text
      >
        <span class="mr-2">Latest Release</span>
        <v-icon>mdi-open-in-new</v-icon>
      </v-btn>
      --> 