(function () {
    var app = angular.module('lemonade', ['ngRoute', 'ipCookie', 'ngAnimate']);
    var baseUrl;
    if (window.location.port != "") {
        baseUrl = location.protocol + "//" + window.location.hostname + ":" + window.location.port + "/api/v1"
    } else {
        baseUrl = location.protocol + "//" + window.location.hostname + "/api/v1"
    }
    var basePath;
    if (window.location.port != "") {
        basePath = location.protocol + "//" + window.location.hostname + ":" + window.location.port;
    } else {
        basePath = location.protocol + "//" + window.location.hostname;
    }
    app.run(function ($rootScope) {
        $rootScope.loading = false;
        $rootScope.httpCount = 0;
        $rootScope.bad = 0;
    });
    // adding the interceptor for the session validation
    app.factory('myHttpInterceptor', function ($q, $window, $rootScope, $location) {
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
                        $location.path("/login");
                }
                return $q.reject(rejection);
            }
        };
    });
    // facebook comments directive
    app.directive('dynFbCommentBox', function () {
        function createHTML(href, numposts, colorscheme) {
            return '<div class="fb-comments" ' +
                'data-href="' + href + '" ' +
                'data-numposts="' + numposts + '" ' +
                'data-colorsheme="' + colorscheme + '">' +
                '</div>';
        }


        return {
            restrict: 'A',
            scope: {},
            link: function postLink(scope, elem, attrs) {
                attrs.$observe('pageHref', function (newValue) {
                    var href = newValue;
                    var numposts = attrs.numposts || 5;
                    var colorscheme = attrs.colorscheme || 'light';
                    elem.html(createHTML(href, numposts, colorscheme));

                    console.log("testing");
                    FB.XFBML.parse(elem[0]);

                });
            }
        };
    });
    app.directive('dyFbLike', function () {
        function createHTML(href) {
            return '<div class="fb-share-button" data-href="' + href + '" data-layout="button_count"></div>';
        }

        return {
            restrict: 'A',
            scope: {},
            link: function postLink(scope, elem, attrs) {
                attrs.$observe('pageHref', function (newValue) {
                    var href = newValue;
                    elem.html(createHTML(href));
                    FB.XFBML.parse(elem[0]);
                });
            }
        };
    });
    app.directive('dyTwShare', function () {
        function createHTML(href) {
            return ' <a href="https://twitter.com/share" class="twitter-share-button" data-size="small" data-text="Make groups and Get huge discounts on mobiles #mobiledevices #smartphones" data-url="' + href + '">Tweet</a>';
        }

        return {
            restrict: 'A',
            scope: {},
            link: function postLink(scope, elem, attrs) {
                attrs.$observe('pageHref', function (newValue) {
                    var href = newValue;
                    elem.html(createHTML(href));
                });
            }
        };
    });
    app.directive('dyGplusShare', function () {
        function createHTML(href) {
            return '<div class="g-plus" data-action="share" data-annotation="bubble" data-height="20" data-href="' + href + '"></div>';
        }

        return {
            restrict: 'A',
            scope: {},
            link: function postLink(scope, elem, attrs) {
                attrs.$observe('pageHref', function (newValue) {
                    var href = newValue;
                    elem.html(createHTML(href));
                });
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
            when('/login', {
                templateUrl: 'public/partials/login.html',
                controller: 'LoginPageController'
            }).
            //when('/share/:dealId', {
            //    templateUrl: 'public/partials/share.html',
            //    controller: 'SharePageController'
            //}).
            when('/share-widget/:dealId', {
                templateUrl: 'public/partials/shareWidget.html',
                controller: 'ShareWidgetController'
            }).
            when('/sign-up/success', {
                templateUrl: 'public/partials/signUpSuccess.html',
                controller: 'SignUpSuccessPageController'
            }).
            when('/privacy', {
                templateUrl: 'public/partials/privacyPolicy.html',
                controller: 'PrivacyPolicyPageController'
            }).
            //when('/dashboard', {
            //    templateUrl: 'public/partials/dashboard.html',
            //    controller: 'DashboardPageController'
            //}).
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
        $scope.login = function () {
            $location.path("/login");
        };
        $scope.policy = function () {
            $location.path("/privacy");
        };
    });

    app.controller('PrivacyPolicyPageController', function ($scope, $http, $location, $window) {

        $scope.init = function(){
            $window.scrollTo(0, 0);
        };
        $scope.goToDashboard = function () {
            $location.path("/");
        };
    });

    app.controller('LoginPageController', function ($scope, $http, $location, $window, ipCookie) {
        $scope.loginStatus = {};
        $scope.user = {};
        $scope.goToDashboard = function () {
            $window.location.href = '/';
        };
        $scope.login = function () {
            var btn = $('#loginButton').button('loading');
            $http.post(baseUrl + '/user/login', $scope.user).success(function (data, status) {
                //console.log(data);
                btn.button('reset');
                if (data.success) {
                    $scope.loginStatus.is_success = true;
                    ipCookie("lemonades_session_key", data.user.session_key, {
                        expires: 1,
                        path: '/'
                    });
                    $window.location.href='/dashboard';
                    return;
                }
                $scope.loginStatus.message = data.message;
                $scope.loginStatus.is_success = false;

            }).error(function (data, status) {
                //console.log(data);
                $scope.loginStatus.message = data.message;
                $scope.loginStatus.is_success = false;
                btn.button('reset');
            });
        };
    });

    app.controller('DashboardPageController', function ($scope, $http, $location, $window, ipCookie, $interval, $rootScope) {
        $scope.deal = {};
        $scope.contentLoaded = false;

        $scope.logout = function () {
            $http.post(baseUrl + '/user/logout', null).success(function (data, status) {
                //console.log(data);
                if (data.success) {
                    ipCookie.remove("lemonades_session_key");
                    $window.location.href='/#!/login';
                }
            }).error(function (data, status) {
                //console.log(data);
            });
        };

        $scope.init = function () {
            $.getScript('http://platform.twitter.com/widgets.js');
            $http.get(baseUrl + '/user/deals').success(function (data, status) {
                //console.log(data);
                if (data.success) {
                    $scope.deal = data.deal;
                    $scope.contentLoaded = true;
                }
            }).error(function (data, status) {
                //console.log(data);
            });
        };
    });

    app.controller('SharePageController', function ($scope, $http, $location, $window, $routeParams) {
        $scope.dealId = $routeParams.dealId;
        $scope.deal = {};

        $scope.goToRegister = function () {
            $window.location.href = "/#!/sign-up";
        };

        $scope.goToDashboard = function () {
            $window.location.href = '/';
        };

        $scope.init = function (dealId) {
            $window.scrollTo(0, 0);
            $http.get(baseUrl + '/deal/' + dealId).success(function (data, status) {
                //console.log(data);
                if (data.success) {
                    $scope.deal = data.deal;
                }
            }).error(function (data, status) {
                //console.log(data);
            });
        };
    });

    app.controller('ShareWidgetController', function ($scope, $http, $location, $window, $routeParams) {
        $scope.dealId = $routeParams.dealId;
    });

    app.controller('SignUpSuccessPageController', function ($scope, $http, $location, $window) {
        $scope.init = function () {
            $window.scrollTo(0, 0);
        };
        $scope.goToDashboard = function () {
            $window.location.href = '/';
        };
    });

    app.controller('SignUpPageController', function ($scope, $http, $routeParams, $window, $location) {
        $scope.signUpStatus = {};
        $scope.user = {city: "Pune"};

        $scope.init = function () {
            $window.scrollTo(0, 0);
        };

        $scope.policy = function () {
            $location.path("/privacy");
        };

        $scope.goToDashboard = function () {
            $window.location.href = '/';
        };

        $scope.signUp = function () {
            var btn = $('#signUpButton').button('loading');
            $scope.signUpStatus.is_tried = true;
            $http.post(baseUrl + '/user', $scope.user).success(function (data, status) {
                //console.log(data);
                btn.button('reset');
                if (data.success) {
                    $scope.signUpStatus.is_success = true;
                    $location.path("/sign-up/success");
                    return;
                }
                $scope.signUpStatus.message = data.message;
                $scope.signUpStatus.is_success = false;

            }).error(function (data, status) {
                //console.log(data);
                $scope.signUpStatus.message = data.message;
                $scope.signUpStatus.is_success = false;
                btn.button('reset');
            });
        };
    });

})();