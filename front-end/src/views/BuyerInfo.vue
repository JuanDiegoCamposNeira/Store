<!---------------------------------------------->
<!--              Template                    --> 
<!---------------------------------------------->
<template>
    <v-container fluid>
        <v-row class="justify-center">

            <!-- Buyer's info -->
            <v-col class="buyer-info mb-4 d-flex justify-space-around" 
                   cols="11" >
                <p class="font-weight-light display-1"> 
                    Comprador : 
                    <span class="font-weight-normal display-2">
                        {{ buyerInfo.name }} 
                    </span> 
                </p> 
                <p class="font-weight-light display-1"> 
                    Edad:
                    <span class="font-weight-normal display-2">
                         {{ buyerInfo.age }} 
                    </span>
                </p>
            </v-col>            

            <!-- Suggestions -->
            <v-col class="suggestions mb-4 d-flex flex-flow-vertical flex-wrap"  
                   cols="11" >
                <!-- Title --> 
                <p class="display-1 font-weight-bold"> Sugerencias </p>
                <!-- Suggestions --> 
                <Suggestion class="suggestion" v-for="suggestion in suggestionsNames" 
                            :key="suggestion.uid" 
                            :productSuggestion="suggestion" 
                            :productsSuggested="suggestions[suggestion]" />
            </v-col>  

            <!-- Transactions history -->
            <v-col style="border-right: 3px solid #190862; " class="history mb-4" 
                   cols="11" 
                   xs="11"
                   sm="11" 
                   md="7"
                   lg="7"
                   xl="7" >
               <!-- Title -->  
               <p class="text-center font-weight-bold display-1"> Historial </p>
               <!-- History -->
               <div v-for="transaction in history" 
                    :key="transaction.id"
                    class="d-flex justify-space-around align-center mb-2 transaction-wrapper">
                    <span>
                        <span class="font-weight-bold"> Fecha : </span> {{ transaction.Date }} 
                    </span>
                    <span>
                        <span class="font-weight-bold"> Cantidad Productos : </span> {{ transaction.products.length }} 
                    </span>
                    <span>
                        <span class="font-weight-bold"> Precio total : </span> {{ transaction.totalPrice }}
                    </span>
               </div>
            </v-col>

            <!-- Same Ip -->
            <v-col style="border-left: 3px solid #190862; " class="sameIp mb-4" 
                   cols="11"
                   xs="11"
                   sm="11" 
                   md="4" >
               <!-- Title -->  
               <p class="text-center font-weight-bold display-1"> Misma Ip </p>
                <!-- Same Ip buyers -->
               <div v-for="transaction in sameIp" 
                    :key="transaction.id"
                    class="d-flex justify-space-around align-center mb-2 transaction-wrapper">
                    <span>
                        <span class="font-weight-bold justify-start"> Nombre: </span> {{ transaction.buyer.name }} 
                    </span>
                    <span>
                        <span class="font-weight-bold"> Id: </span> {{ transaction.buyer.uid }}
                    </span>
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
import Suggestion from "../components/Suggestion.vue";

export default {
    name: "BuyerInfo", 

    methods: () => {
        return {}
    },

    mounted() {
        // Make request to backend to get the buyer information 
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
                this.suggestionsNames = Object.keys(this.suggestions); 
                // Save personal information about the user
                this.buyerInfo.name = response.data.history[0].buyer.name;
                this.buyerInfo.age = response.data.history[0].buyer.age;
                // Debug 
                console.log(this.suggestions);
             });
    },

    components: {
        Suggestion, 
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
            suggestionsNames: [], 
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
    .buyer-info, .suggestions, .history, .sameIp {
        background-color: #F3F6FF !important; 
        border-radius: 10px; 
    }
    .suggestion {
        width: 400px;
    }
    .transaction-wrapper {
        background-color: #4800FF !important; 
        color: #ffffff !important; 
        border-radius: 10px;
        min-height: 60px;
    }
</style>