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

        <!-- Buyers info -->
        <buyer v-for="buyer in this.buyersData" 
                :key="buyer.name" 
                :name="buyer.name" :age="buyer.age" :uid="buyer.uid" />

        <!-- View when there are no buyers --> 
        <div v-if="!this.buyersData || this.buyersData.length == 0">
            <p class="display-1 white--text mt-10 text-center font-weight-light"> 
                Aún no hay compradores en el sistema
            </p>
        </div>

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
            // Request buyers information
            axios.get("http://localhost:3000/buyers")
                .then( response => {
                    // Parse information to display for all buyers
                    let buyers = []
                    response.data.buyers.forEach(person => {
                        buyers.push({name: person.name, age: person.age, uid: person.uid })
                    });
                    // Save buyers data 
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