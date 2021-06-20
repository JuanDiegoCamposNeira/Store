<template>
    <div> 
        <v-container> 
            <v-row> 
                <!-- Title --> 
                <v-col>
                    <h2> {{ title }} </h2>
                </v-col>

                <!-- File input -->
                <v-col>
                    <v-file-input
                        label="Seleccione el archivo con extensión .json"
                        @change="handleAddFile"
                    ></v-file-input>
                </v-col> 

                <!-- Submit button --> 
                <v-col> 
                    <v-btn @click="handleAddFile"> Agregar </v-btn>
                </v-col>

            </v-row>
        </v-container>

    </div>
</template>

<!---------------------------------------------->
<!--               Script                     --> 
<!---------------------------------------------->
<script>

    export default {
        name: 'Add',

        mounted() {
            console.log("Add mounted")
        },

        props: {
            title: { type: String, requred: true },
        },

        data: () => ({
            test: null
        }),

        methods: {
            handleAddFile: (file) => {
                console.log(file)
                const fileReader = new FileReader()
                console.log("File reader result", fileReader.result) 
            }
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
</style>