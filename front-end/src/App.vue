<template>
  <v-app>

    <!---------------------------------------------->
    <!--               Navbar                     --> 
    <!----------------------------------------------> 
    <v-app-bar app dark>
      <!-- Title --> 
      <v-toolbar-title> <span> <span class="T">T</span> ienda </span> </v-toolbar-title>
      <v-spacer></v-spacer>

      <!-- Navigation links --> 
      <router-link to="/add"> Añadir </router-link>
      <v-spacer></v-spacer>
      <router-link to="/buyers"> Compradores </router-link>

    </v-app-bar>

    <!---------------------------------------------->
    <!--               Router views               --> 
    <!---------------------------------------------->
    <v-main>
      <router-view/>
    </v-main>

    <!---------------------------------------------->
    <!--                 Footer                   --> <!-- Choose one --> 
    <!---------------------------------------------->
    <v-footer padless>
      <v-col class="text-center" cols="12" >
        Made with ♥ by JD 
    </v-col>
    </v-footer>
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
  .v-app-bar {
    background-color: #190862 !important;
  }
  .T {
    font-family: 'Lobster', cursive;
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