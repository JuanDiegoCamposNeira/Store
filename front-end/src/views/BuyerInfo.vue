<!---------------------------------------------->
<!--              Template                    --> 
<!---------------------------------------------->
<template>
    <v-container fluid>
        <v-row>

            <!-- Buyer's info -->
            <v-col class="debug buyerInfo" 
                   cols="10" >
                <p> Comprador : {{ buyerInfo.name }} </p>
            </v-col>            

            <!-- Suggestions -->
            <v-col class="debug suggestions d-flex flex-flow-horizontal" 
                   cols="11" >
                <!-- Suggestions --> 
                <p v-for="suggestion in suggestions" :key="suggestion.uid"> 
                    || {{ 1 }} || - 
                </p>
            </v-col>  

            <!-- Transactions history -->
            <v-col class="debug history" 
                   sm="11" 
                   md="6" 
                   lg="5" >
               <!-- Title -->  
               <p class="text-center"> Historial </p>
               <!-- History -->
               <div v-for="transaction in history" 
                    :key="transaction.id">
                Id : {{ transaction.uid }} | 
                Cantidad Productos : {{ transaction.products.length }} | 
                Precio total : {{ transaction.totalPrice }}
               </div>
            </v-col>

            <!-- Same Ip -->
            <v-col class="debug sameIp" 
                   sm="11" 
                   md="5" 
                   lg="5" >
               <!-- Title -->  
               <p> Misma Ip </p>
                <!-- Same Ip buyers -->
               <div v-for="transaction in sameIp" :key="transaction.id">
                   Nombre: {{ transaction.buyer.name }} | 
                   Id: {{ transaction.buyer.uid }}
               </div>
            </v-col>
        </v-row>
    </v-container>
</template>

<!---------------------------------------------->
<!--               Script                     --> 
<!---------------------------------------------->
<script>
import axios from "axios";  

export default {
    name: "BuyerInfo", 

    methods: () => {
        return {}
    },

    mounted() {
        console.log("Single buyer mounted") 
        console.log(this.$route.params.buyerId)
        axios.get(`http://localhost:3000/buyer/${this.$route.params.buyerId}`)
             .then( response => {
                // Add total price to all transactions 
                response.data.history.forEach( transaction => {
                    let totalCost = 0;
                    transaction.products.forEach( product => totalCost += parseInt(product.price));
                    transaction.totalPrice = totalCost; 
                });

                // Save information
                this.sameIp = response.data.sameIp;
                this.history = response.data.history;
                this.suggestions = response.data.Suggestions;
                // Save personal information about the user
                this.buyerInfo.name = response.data.history[0].buyer.name;
                this.buyerInfo.age = response.data.history[0].buyer.age;
                // Debug 
                console.log(this.suggestions);
             });
    },

    components: {
    },

    data: () => {
        return { 
            buyerInfo: {
                name: "", 
                age: "",
            },
            sameIp: [], 
            history: [],
            suggestions: [],
        } 
    }
}
</script>

<!---------------------------------------------->
<!--                Style                     --> 
<!---------------------------------------------->
<style scoped>
    .debug {
        border: 1px solid black; 
    }
    .buyerInfo, .suggestions, .history, .sameIp {
        background-color: #F3F6FF !important; 
        border-radius: 10px; 
    }
</style>