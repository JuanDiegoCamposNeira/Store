<template>
    <v-container fluid class="suggestion-wrapper mb-4 pa-0">
        
        <!-- Title --> 
        <v-row class="product-suggestion">
            <v-col cols="4">
                Por ordenar: 
            </v-col>
            <v-col cols="8" class="name">
                {{ productSuggestion.length > 29 ? productSuggestion.substring(0, 29) + "..." : productSuggestion }}
            </v-col>
        </v-row>

        <!-- Products -->
        <div v-for="suggestion in productsSuggested" :key="suggestion.uid"> 
            <v-container fluid>
                <v-row>
                    <v-col cols="8" class="product-info"> 
                        <p class="label"> Producto </p> 
                        <p> {{ suggestion.name.length > 32 ? suggestion.name.substring(0, 32) + "..." : suggestion.name }} </p> 
                    </v-col> 
                    <v-col cols="1" class="separator"> | </v-col>
                    <v-col cols="" class="product-info"> 
                        <p class="label">Precio</p> 
                        <p>{{ suggestion.price }}</p> 
                    </v-col>
                </v-row>
            </v-container>
        </div>

    </v-container>
</template>


<script>
    export default {
        name: "Suggestion", 
        props: {
            productSuggestion: String, 
            productsSuggested: Array, 
        }
    }
</script>


<style scoped>
    .label {
        font-size: 0.8rem;
        font-weight: 350;
        padding-top: 0px; 
        margin-bottom: -2px;
    }
    .product-info {
        color: #4800FF; 
        font-weight: 700;
        font-size: 0.85rem;
        padding-top: 0px;
        padding-bottom: 0px; 
        margin-bottom: -4px;
    }
    .separator {
        color: #4800FF; 
        font-size: 30px;
        font-weight: 295;
        padding-top: 0px; 
        padding-bottom: 0px; 
    }
    .product-suggestion {
        font-size: 0.9rem !important;
        text-align: center !important;
        margin: -2px auto -3px auto; 
    }
    .product-suggestion .name {
        font-weight: bold;
    }
    .suggestion-wrapper {
        border: 1.5px solid #4800FF; 
        border-radius: 10px; 
    }
</style>