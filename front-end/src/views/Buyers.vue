<!---------------------------------------------->
<!--              Template                    --> 
<!---------------------------------------------->
<template>
    <div>
        <!-- Title --> 
        <p class="title d-flex justify-center align-center"> Compradores </p>
        <!-- Buyers --> 
        <buyer v-for="buyer in this.buyersData" 
                :key="buyer.name" 
                :name="buyer.name" :age="buyer.age" :uid="buyer.uid" />
    </div>
</template>

<!---------------------------------------------->
<!--               Script                     --> 
<!---------------------------------------------->
<script>

import axios from 'axios';
import Buyer from '../components/Buyer.vue';

export default {
    name: "Buyers", 

    mounted() {
        console.log("Buyers mounted") 
        axios.get("http://localhost:3000/buyers")
             .then( response => {
                 // Traverse response data and create person object
                 let buyers = []
                 response.data.buyers.forEach(person => {
                     buyers.push({name: person.name, age: person.age, uid: person.uid })
                 });
                 this.buyersData = buyers
             })
    },

    components: {
        Buyer,
    }, 

    data: () => ({
        buyersData: [],
    })
}
</script>

<!---------------------------------------------->
<!--                Style                     --> 
<!---------------------------------------------->
<style scoped>
.title {
    width: 100%;
    height: 4rem;
    /* Text */
    color: #F3F6FF !important; 
    font-size: 3.5rem !important;
    font-weight: 350; 
}
</style>