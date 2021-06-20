<!---------------------------------------------->
<!--              Template                    --> 
<!---------------------------------------------->
<template>
    <div>
        <h1> Hi {{ buyerInfo.name }} :) </h1>
        <h2> SameIp </h2>
        <p> {{ sameIp }} </p>
        <h2> History </h2>
        <p> {{ history }} </p>
        <h2> Suggestions </h2>
        <p> {{ suggestions }} </p> 
    </div>
</template>

<!---------------------------------------------->
<!--               Script                     --> 
<!---------------------------------------------->
<script>
import axios from "axios";  

export default {
    name: "Buyers", 

    mounted() {
        console.log("Single buyer mounted") 
        console.log(this.$route.params.buyerId)
        axios.get(`http://localhost:3000/buyer/${this.$route.params.buyerId}`)
             .then( response => {
                // Save information
                this.sameIp = response.data.sameIp
                this.history = response.data.history
                this.suggestions = response.data.Suggestions
                // Save personal information about the user
                this.buyerInfo.name = response.data.history[0].buyer.name
                this.buyerInfo.age = response.data.history[0].buyer.age
             });
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