Feature: Weather API
  In order to know weather of a specific location
  As an API user
  I need to be able to request weather

  Scenario: Get weather by city name
    Given I send GET request to "http://api.openweathermap.org/data/2.5/weather?q=London&appid=YOUR_API_KEY"
    Then The response status code should be 200
    And the response body should contain the following JSON properties:
      | coord.lon             | -0.1257       |
      | coord.lat             | 51.5085       |
      | weather.0.id          | 803           |
      | weather.0.main        | Clouds        |
      | weather.0.description | broken clouds |
      | main.temp             | 293.45        |
      | main.feels_like       | 293.15        |
      | main.temp_min         | 291.96        |
      | main.temp_max         | 294.77        |
      | main.pressure         | 1005          |
      | main.humidity         | 62            |
      | visibility            | 10000         |
      | wind.speed            | 7.2           |
      | wind.deg              | 300           |
      | clouds.all            | 75            |
      | dt                    | 1688211936    |
      | sys.type              | 2             |
      | sys.id                | 2075535       |
      | sys.country           | GB            |
      | sys.sunrise           | 1688183234    |
      | sys.sunset            | 1688242862    |
      | timezone              | 3600          |
      | id                    | 2643743       |
      | name                  | London        |
      | cod                   | 200           |

  Scenario: Get weather by lat and long
    Given I send GET request to "http://api.openweathermap.org/data/2.5/weather?lat=44.34&lon=10&APPID=YOUR_API_KEY"
    Then The response status code should be 200
    And the response body should contain the following JSON properties:
      | coord.lon             | 10            |
      | coord.lat             | 44.34         |
      | weather.0.id          | 803           |
      | weather.0.main        | Clouds        |
      | weather.0.description | broken clouds |
      | base                  | stations      |
      | main.temp             | 293.45        |
      | main.feels_like       | 293.15        |
      | main.temp_min         | 291.96        |
      | main.temp_max         | 294.77        |
      | main.pressure         | 1005          |
      | main.humidity         | 62            |
      | visibility            | 10000         |
      | wind.speed            | 7.2           |
      | wind.deg              | 300           |
      | clouds.all            | 75            |
      | dt                    | 1688211936    |
      | sys.type              | 2             |
      | sys.id                | 2028413       |
      | sys.country           | IT            |
      | sys.sunrise           | 1688183234    |
      | sys.sunset            | 1688242862    |
      | timezone              | 7200          |
      | id                    | 3182497       |
      | name                  | Bagnone       |
      | cod                   | 200           |

  Scenario: Get weather by invalid lat
    Given I send GET request to "http://api.openweathermap.org/data/2.5/weather?lat=700&lon=40&APPID=YOUR_API_KEY"
    Then The response status code should be 400
    And the response should match json:
      """
      {
          "cod": "400",
          "message": "wrong latitude"
      }
      """

  Scenario: Get weather by invalid lon
    Given I send GET request to "http://api.openweathermap.org/data/2.5/weather?lat=30&lon=7000&APPID=YOUR_API_KEY"
    Then The response status code should be 400
    And the response should match json:
      """
      {
          "cod": "400",
          "message": "wrong longitude"
      }
      """

  Scenario: Get weather by invalid city name
    Given I send GET request to "http://api.openweathermap.org/data/2.5/weather?q=hjhjhj&appid=YOUR_API_KEY"
    Then The response status code should be 404
    And the response should match json:
      """
      {
          "cod": "404",
          "message": "city not found"
      }
      """

  Scenario: Get weather by city name in XML format
    Given I send GET request to "http://api.openweathermap.org/data/2.5/weather?q=London&appid=YOUR_API_KEY&mode=xml"
    Then The response status code should be 200
    And the response body should contain the following XML properties:
      | current.city.coord.lat     | 51.5085 |
      | current.city.coord.lon     | -0.1257 |
      | current.city.country       | GB      |
      | current.city.timezone      | 3600    |
      | current.city.id            | 2643743 |
      | current.city.sun.rise      |         |
      | current.city.sun.set       |         |
      | current.temperature.min    |         |
      | current.temperature.max    |         |
      | current.temperature.value  |         |
      | current.temperature.unit   |         |
      | current.feels_like.unit    |         |
      | current.feels_like.value   |         |
      | current.humidity.unit      |         |
      | current.humidity.value     |         |
      | current.pressure.unit      |         |
      | current.pressure.value     |         |
      | current.wind.speed.value   |         |
      | current.wind.speed.unit    |         |
      | current.wind.speed.name    |         |
      | current.clouds.name        |         |
      | current.clouds.value       |         |
      | current.visibility.value   |         |
      | current.precipitation.mode |         |
      | current.weather.value      |         |
      | current.lastupdate.value   |         |
