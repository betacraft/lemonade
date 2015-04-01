(function () {
    var app = angular.module('lemonade', ['ngRoute']);
    app.run(function ($rootScope) {
        $rootScope.loading = false;
        $rootScope.httpCount = 0;
        $rootScope.currentlyShowing = 0;
        $rootScope.userEmail = "";
        $rootScope.userName = "";
        $rootScope.selectedApp = {};
    });
    // adding the interceptor for the session validation
    app.factory('myHttpInterceptor', function ($q, $window, $rootScope) {
        return {
            // optional method
            'request': function (config) {
                console.log(config);
                if (!config.ignoreLoadingFlag) {
                    // do something on success
                    console.log("changin flag");
                    $rootScope.loading = true;
                }
                return config;
            },
            // optional method
            'requestError': function (rejection) {
                $rootScope.loading = false;
                // do something on error
                if (canRecover(rejection)) {
                    return responseOrNewPromise
                }
                return $q.reject(rejection);
            },
            // optional method
            'response': function (response) {
                $rootScope.loading = false;
                // do something on success
                console.log("inside intercepter success");
                return response;
            },
            // optional method
            'responseError': function (rejection) {
                $rootScope.loading = false;
                var status = rejection.status;
                if (status == 401) {
                    console.log($window.location.pathname);
                    if ($window.location.pathname != '/login')
                        $window.location.href = '/login';
                }
                return $q.reject(rejection);
            }
        };
    });
    // adding http interceptor
    app.config(function ($httpProvider) {
        $httpProvider.interceptors.push('myHttpInterceptor');
        $httpProvider.defaults.useXDomain = true;
        delete $httpProvider.defaults.headers.common['X-Requested-With'];
    });
    // route information
    app.config(['$routeProvider', '$locationProvider', function ($routeProvider, $locationProvider) {
        $routeProvider.
            when('/', {
                templateUrl: 'public/partials/landing.html',
                controller: 'LandingPageController'
            }).
            when('/sign-up', {
                templateUrl: 'public/partials/signUp.html',
                controller: 'SignUpPageController'
            }).
            when('/dashboard', {
                templateUrl: 'public/partials/dashboard.html',
                controller: 'DashboardPageController'
            }).
            otherwise({
                redirectTo: '/'
            });
        $locationProvider
            .html5Mode(false)
            .hashPrefix('!');
    }]);


    app.controller('LandingPageController', function ($scope, $http, $location, $window) {
        $scope.signUp = function () {
            $location.path("/sign-up");
        };
    });

    app.controller('DashboardPageController', function ($scope, $http, $location, $window) {
        $scope.signUp = function () {
            $location.path("/sign-up");
        };
    });

    app.controller('SignUpPageController', function ($scope, $http, $routeParams, $window) {
        $scope.signUpStatus = {};
        $scope.signUpStatus.is_tried = false;
        $scope.signUpStatus.is_success = false;
        $scope.user = {address:{city:"Pune"}};

        $scope.goToLogin = function () {
            $window.location.href = '/login';
        };

        $scope.signup = function () {
            $scope.signUpStatus.is_tried = true;
            $http.post(baseUrl + '/user/signup', $scope.user).success(function (data, status) {
                console.log(data);
                if (data.success) {
                    $scope.signUpStatus.is_success = true;
                    return;
                }
                $scope.signUpStatus.message = data.message;
                $scope.signUpStatus.is_success = false;
            }).error(function (data, status) {
                console.log(data);
                $scope.signUpStatus.message = data.message;
                $scope.signUpStatus.is_success = false;
            });
        };
    });

})();