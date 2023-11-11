$(document).ready(() => {
    $('#go').click(() => {
        getWeather($('#city').val())
    })
    $('#city').keypress(function (event) {
        var keycode = (event.keyCode ? event.keyCode : event.which);
        if (keycode == '13') {
            getWeather($('#city').val())
        }
    });

    function getWeather(city) {
        $.ajax({
            type: 'get',
            url: '/weather/' + city,
            success: function (res){
                result = JSON.parse(res)
                $('#result').show()
                $('#weather_icon').attr("src","https://" + result.current.condition.icon)
                $('#weather_text').html(result.current.condition.text)
                $('#city_name').html(result.location.name)
                $('#country_name').html(result.location.country)
                $('#temp').html(result.current.temp_c + "&deg;C&nbsp;-&nbsp;" + result.current.temp_f + "&deg;F")
                $('#feels_like').html(result.current.feelslike_c + "&deg;C&nbsp;-&nbsp;" + result.current.feelslike_f + "&deg;F")
            }
        })
    }
})