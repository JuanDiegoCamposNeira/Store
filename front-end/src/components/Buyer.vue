<template>
  <v-container class="buyer-wrapper mb-2 mt-2 ">
    <v-row class="d-flex align-center justify-center">
      <!-- Name --> 
      <v-col cols="5" class="d-flex align-center"> 
        <span class="title mr-4"> Cliente : </span> 
        <span class="info">{{ name }}</span> 
      </v-col>
      <!-- Age -->
      <v-col cols="3" class="d-flex align-center"> 
        <span class="title mr-4"> Edad : </span> 
        <span class="info">{{ age.toString() }}</span>
      </v-col>
      <!-- See buyer info button --> 
      <v-col cols="3"> 
        <v-btn class="buyerInfoBtn" 
              @click="handleBuyerInfo"> Ver información </v-btn>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
export default {
  name: "Buyer", 

  data() {
    return { 
    }
  },

  props: {
    name: String, 
    age: Number,
    uid: String, 
  }, 

  methods: {
    handleBuyerInfo() {
      this.$router.push({ path: `/buyer/${ this.$props.uid}` })
    }
  }
}
</script>

<style scoped>
.buyer-wrapper {
  width: 90%;
  background-color: #F3F6FF !important;  
  border-radius: 15px;
}
.title {
  font-weight: 350;
  font-size: 30px !important;
}
.info {
  background-color: #F3F6FF !important;
  font-size: 40px;
}
.buyerInfoBtn {
  width: 100% !important;
  font-size: 70% !important;
  background-color: #4800FF !important; 
  color: #ffffff !important; 
}
</style>