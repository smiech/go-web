/* Thanks to CSS Tricks for pointing out this bit of jQuery
    https://css-tricks.com/equal-height-blocks-in-rows/
    It's been modified into a function called at page load and then each time the page is resized. One large modification was to remove the set height before each new calculation. */

equalheight = function(container) {
    var currentTallest = 0,
        currentRowStart = 0,
        rowDivs = new Array(),
        $el,
        topPosition = 0;
    $(container).each(function() {
        $el = $(this);
        $($el).height("auto");
        topPostion = $el.position().top;

        if (currentRowStart != topPostion) {
            for (currentDiv = 0; currentDiv < rowDivs.length; currentDiv++) {
                rowDivs[currentDiv].height(currentTallest);
            }
            rowDivs.length = 0; // empty the array
            currentRowStart = topPostion;
            currentTallest = $el.height();
            rowDivs.push($el);
        } else {
            rowDivs.push($el);
            currentTallest = currentTallest < $el.height()
                ? $el.height()
                : currentTallest;
        }
        for (currentDiv = 0; currentDiv < rowDivs.length; currentDiv++) {
            rowDivs[currentDiv].height(currentTallest);
        }
    });
};

isMobile = function() {
    if (window.matchMedia) {
        return window.matchMedia('(max-width: 769px)').matches;
    }
};

$(document).ready(function() {
// Show modals
    try {
        $('.jsModal').modaal({
            fullscreen : isMobile( $(this) ),
            // custom_class :
        });

        $('.jsModalBlue').modaal({
            fullscreen : isMobile( $(this) ),
            custom_class : 'modaal--blue',
            before_open : function(e) {
                console.log( $('.high-seventy-two').attr('src', $('.high-seventy-two').attr('src').replace(/\?.*$/,"")+"?x="+Math.random() ) );
            }
        });

        $('.jsModalSlider').modaal({
            fullscreen : isMobile(),
            custom_class : 'has_slider is_onboarding',
            before_open : function(event) {
            },
            after_open : function() {
                var $slider = $('.modaal-container .pm-modal-slider');
                $slider.slick({
                    // Slick config goes here...
                    initialSlide : 0,
                    dots : true,
                    infinite : false,
                    draggable : false,
                    arrows : false,
                    adaptiveHeight : isMobile(),
                    appendDots : $('.pm-modal-slider__footer')
                });
                if ( !isMobile() ) {
                    equalheight(".modaal-container .pm-modal-slider .pm-modal-slide");
                }

                // Controls
                // Next slide button
                $('.modaal-container').find('.jsSliderNext').click(function() {
                    $slider.slick('slickNext');
                });

                // Prev slide button
                $('.modaal-container').find('.jsSliderPrev').click(function() {
                    $slider.slick('slickPrev');
                });
                // Close button
                $('.modaal-container').find('.jsSliderClose').click(function() {
                    $('.jsModalSlider').modaal('close');
                });

            },
            before_close : function() {
                // Destroy Slick as we close.
                $('.modaal-container .pm-modal-slider.slick-initialized').slick('unslick');
            }
        });

        $('.jsModalConfirm').modaal({
            fullscreen : isMobile( $(this) ),
            is_locked : true,
            hide_close : true,
            width : 520,
            custom_class : 'pm-confirm',
            after_open : function() {
                $('.modaal-container').find('.jsConfirmOk').click(function() {
                    console.log('you have confirmed this action');
                    $('.jsModalConfirm').modaal('close');
                });
                $('.modaal-container').find('.jsConfirmCancel').click(function() {
                    console.log('you have canceled this action');
                    $('.jsModalConfirm').modaal('close');
                });
            }
        });

    } catch(err) {
        console.log('The library "modaal.js" is not included on this page.');
        console.log( err );
    }

// Reset HighFive
    $('body').on('click', '.high-seventy-two', function() {
        console.log( $(this).attr('src', $(this).attr('src').replace(/\?.*$/,"")+"?x="+Math.random() ) );
    });

    $('.pm-client').on('click', '.hover-menu-item', function() {
        if (isMobile()) {
            $(this).parents('.hover-menu').toggleClass('is_focused');
        }
    });

// Toggle Password
    $('.pm-field--password').on('change', '.jsPasswordToggle', function() {
        var $this = $(this),
            parent = $this.closest(".pm-field"),
            field = parent.find(".jsPasswordField");

        if ($this.is(':checked')) {
            field.attr('type', 'text');
        } else {
            field.attr('type', 'password');
        }
    });

    // On before slide change
    $('.pm-modal-slider').on('beforeChange', function(event, slick, currentSlide, nextSlide){
        // console.log( $(this).parents('.modaal-content-container').scrollTop( 0 ) );
        $(this).parents('.modaal-content-container').animate({ scrollTop: 0 }, 290, "linear")
    });

    $('.circle-loader').toggleClass('load-complete');
    $('.checkmark').toggle();

    $('body').on('click', '.jsToggle', function() {
        var $this = $(this),
            thisValueA = $this.data('toggle-text'),
            thisValueB = $this.html(),
            $target =  $( $this.data('toggle-target') );

        $target.slideToggle();
        $this.html( thisValueA ).data( 'toggle-text', thisValueB )
    });

});





