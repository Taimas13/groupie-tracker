ymaps.ready(init);

function init() {
    var newMap = new ymaps.Map("map",{
        center: [30, 10],
        zoom: 1
    });
    
    var cities = document.getElementsByClassName('city');

     for (var city of cities) {
        var newCity = city.innerHTML.replace(/ - /g, "-")
        newCity = newCity.replace(/-/g, " ")
        
        var newGeocoder = ymaps.geocode(newCity, {results: 1, prefLang: "en"});

        newGeocoder.then(function(searchResult) {
            newMap.geoObjects.add(searchResult.geoObjects);
        });
    }
}