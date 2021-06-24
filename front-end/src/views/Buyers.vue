<!---------------------------------------------->
<!--              Template                    --> 
<!---------------------------------------------->
<template>
    <div>
        <!-- Title --> 
        <p class="display-3 white--text d-flex justify-center align-center"> Compradores </p>
        <!-- Loading bar --> 
        <div class="d-flex justify-center">
            <v-progress-circular
            class="mt-15"
            indeterminate
            size="80"
            width="10"
            color="primary"
            id="loading-buyers"
            ></v-progress-circular>
        </div>
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
             .finally( () => {
                document.getElementById("loading-buyers").style.display = "none"; 
             });
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
</style>