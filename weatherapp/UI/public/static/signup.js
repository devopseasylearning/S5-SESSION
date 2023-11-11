let password = document.getElementById("password")
    , confirm_password = document.getElementById("confirm_password");

function validatePassword() {
    if (password.value != confirm_password.value) {
        confirm_password.setCustomValidity("Passwords Don't Match");
    } else {
        confirm_password.setCustomValidity('');
    }
}

password.onchange = validatePassword;
confirm_password.onkeyup = validatePassword;
$(document).ready(function (){
    if (location.search.substring(1).split("=")[1] == "userexists"){
        $("#signupresult").show()
        $("#signupresult").html("<strong>Username already exists</strong>")
    }
})