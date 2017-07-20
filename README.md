# Library-Organizer-2.0

# REST API
* Books `/books`
  * Get Books `GET /books`
  * Save Book `PUT /books`
  * Add Book `POST /books`
  * Delete Book `DELETE /books`
  * Export Books `GET /books/books`
  * Export Authors `GET /books/authors`
  * Import Books `POST /books/books`
* Information `/information`
  * Get Statistics `GET /information/statistics`
  * Get Dimensions `GET /information/dimensions`
  * Get Publishers `GET /information/publishers`
  * Get Cities `GET /information/cities`
  * Get States `GET /information/states`
  * Get Countries `GET /information/countries`
  * Get Formats `GET /information/formats`
  * Get Roles `GET /information/roles`
  * Get Series `GET /information/series`
  * Get Languages `GET /information/languages`
  * Get Deweys `GET /information/deweys`
* Libraries `/libraries`
  * Get Libraries `GET /libaries`
  * Get Owned Libraries `GET /libraries/owned`
  * Save Owned Libraries `PUT /libraries/owned`
  * Get Cases `GET /libraries/:libraryid/cases`
  * Save Cases `PUT /libraries/:libraryid/cases`
  * Get Breaks `GET /libraries/breaks`
  * Add Break `POST /libraries/breaks`
  * Save Break `PUT /libraries/breaks`
  * Delete Break `DELETE /libraries/breaks`
* Settings `/settings`
  * Get Settings `GET /settings`
  * Save Settings `PUT /settings`
  * Get Setting `GET /settings/:setting`
* Users `/users`
  * Get Users `GET /users`
  * Login `PUT /users`
  * Register `POST /users`
  * Logout `DELETE /users`
  * Send Password Reset `PUT /users/password`
  * Finish Password Reset `GET /users/password/:token`
  * Get Username `GET /users/username`
