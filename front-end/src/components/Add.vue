<!---------------------------------------------->
<!--              Template                    --> 
<!---------------------------------------------->
<template>
    <div> 

        <v-container> 
            <v-row> 
                <!------------------------------------->
                <!--              Title              --> 
                <!------------------------------------->
                <v-col cols="4" md="4" sm="11" xs="11">
                    <p class="display-1 white--text text-center"> {{ name }} </p>
                </v-col>

                <!------------------------------------->
                <!--        Input information        -->
                <!------------------------------------->
                <v-col class="text-center">
                    <textarea  @change="handleInputChange"
                               name="input" 
                               :id="name"
                               rows="8" 
                               cols="45" 
                               placeholder="Ingresa la información en formato JSON" /> 
                </v-col> 

                <!------------------------------------->
                <!--          Submit button          --> 
                <!------------------------------------->
                <v-col class="d-flex flex-column"> 
                    <v-btn @click="handleAddInfo" :id="`${name}-submit-btn`"> Agregar </v-btn>
                    <v-alert
                        :id="`${name}-alert`"
                        type="success"
                    > {{ responseMessage }} </v-alert>
                </v-col>

            </v-row>
        </v-container>

    </div>
</template>

<!---------------------------------------------->
<!--               Script                     --> 
<!---------------------------------------------->
<script>

    import axios from "axios"; 

    export default {
        name: 'Add',

        props: {
            name: { type: String, required: true },
        },

        data: () => ({
            responseMessage: "", 
            name: "", 
            inputInfo: "",
        }),

        mounted: function() {
            const alertName = this.$props.name + "-alert"; 
            document.getElementById(alertName).style.display = "none"; 
        },

        methods: {

            handleInputChange: function() {
                this.inputInfo = document.getElementById(this.$props.name).value; 
            },

            handleAddInfo: function() {

                // Get the reference to the elements
                let errorAlert = document.getElementById("alert");                  // Error alert
                const successAlertName = this.name + "-alert";                      // Success alert name
                let successAlert = document.getElementById(successAlertName)        // Success alert element
                let selectedDate = document.getElementById("selected-date").value   // Date selected by the user

                // If there is nothing in the textarea, return 
                if (!this.inputInfo || this.inputInfo === "") {
                    errorAlert.style.display = "inline"; 
                    setTimeout( () => { errorAlert.style.display = "none"; }, 5000); 
                    return;  
                }

                // Display progess bar
                document.getElementById("progress").style.display = "inline"

                // Disable submit buttons
                document.getElementById("Productos-submit-btn").disabled = true; 
                document.getElementById("Compradores-submit-btn").disabled = true; 
                document.getElementById("Transacciones-submit-btn").disabled = true; 
                document.getElementById("Productos").disabled = true; 
                // Disable input areas 
                document.getElementById("Compradores").disabled = true; 
                document.getElementById("Transacciones").disabled = true; 

                // Get name 
                let apiEndpoint = this.name
                // Translate apiEndpoint from spanish to english
                if (apiEndpoint === "Productos") apiEndpoint = "products"; 
                else if (apiEndpoint === "Compradores") apiEndpoint = "buyers"; 
                else apiEndpoint = "transactions"; 

                // Check if there is a selected date
                const date = (selectedDate !== "") ? `?date=${selectedDate}` : ''; 
                
                // Make post request to get the buyer info
                axios.post(`http://localhost:3000/${ apiEndpoint + date}`, this.inputInfo)
                     .then( response => { 
                         this.responseMessage = response.data; 
                     })
                     .finally( err => {
                         // DIsplay error if exists
                         if (err) console.log(err); 

                         // Clear input
                         document.getElementById(this.name).value = "";  
                         this.inputInfo = ""; 

                         // Enable submit buttons
                         document.getElementById("Productos-submit-btn").disabled = false; 
                         document.getElementById("Compradores-submit-btn").disabled = false; 
                         document.getElementById("Transacciones-submit-btn").disabled = false;                         
                         // Enable input areas
                         document.getElementById("Productos").disabled = false; 
                         document.getElementById("Compradores").disabled = false; 
                         document.getElementById("Transacciones").disabled = false; 

                         // Hide progress bar
                         document.getElementById("progress").style.display = "none"; 
                         
                         // Show success message
                         successAlert.style.display = "inline"; 
                         setTimeout( () => {
                             successAlert.style.display = "none";
                         }, 3000); 
                     });
            }, 
        }
    };

</script>

<!---------------------------------------------->
<!--                Style                     --> 
<!---------------------------------------------->
<style scoped>
    .v-btn {
        background-color: #4800FF !important;
        color: #ffffff !important; 
    }
    textarea {
        border: 1px solid #ffffff;
        border-radius: 5px;  
        color: #ffffff; 
    }

</style>