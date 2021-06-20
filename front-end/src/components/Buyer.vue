<template>
  <div>
      <h3> Name: {{ name }} ||| Age: {{ age.toString() }}  </h3>
      <v-btn @click="handleBuyerInfo"> Ver información </v-btn>
  </div>
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