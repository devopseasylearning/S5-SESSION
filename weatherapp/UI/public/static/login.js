$(document).ready(function (){
    if (location.search.substring(1).split("=")[1] == "invalidcreds"){
        $("#loginresult").show()
        $("#loginresult").html("<strong>Invalid credentials</strong>")
    }
})
