var rubberDoc = (function () {
    var eventNames = {
        tabShown: 'RubberDocTabShown',
        tabAlreadyShown: 'RubberDocTabAlreadyShown',

        collapsibleOpened: 'RubberDocCollapsibleOpened',
        collapsibleClosed: 'RubberDocCollapsibleClosed'
    };

    /**
     * needs to have a DOM structure like this:
     *
     * <ul>
     *     <li class="rd-collapsible">
     *          <-- header link can be nested, does not need to be a direct child of the <li> //-->
     *          <a href="#" data-rd-collapsible="link">header 1</a>
     *
     *          <-- content need to be a direct child of the <li> //-->
     *          <div data-rd-collapsible="content">
     *              content to collapse
     *           </div>
     *     </li>
     *
     *     <-- multiple collapsible areas needs to be siblings //-->
     *     <li class="rd-collapsible">
     *          ...
     *     </li>
     * </ul>
     *
     *
     * @param {jQuery} $moduleElement
     * @param {Object} options
     * @constructor
     */
    function CollapsibleManager($moduleElement, options) {
        this.$moduleElement = $moduleElement;
        this.options = $.extend(true, {}, this.defaultOptions, options);
        this.setEvents();
    }

    CollapsibleManager.prototype = {
        defaultOptions: {
            activeClass: 'rd-active',
            collapsibleSelector: '.rd-collapsible'
        },

        setEvents: function() {
            var that = this;

            this.$moduleElement.on('click', '[data-rd-collapsible=link]', function (e) {
                e.preventDefault();
                that.onClick($(this));
            });
        },

        /**
         * @param {jQuery} $link
         */
        onClick: function($link) {
            var $collapsible = $link.closest(this.options.collapsibleSelector);
            if ($collapsible.hasClass(this.options.activeClass)) {
                this.close($collapsible);
            } else {
                this.open($collapsible);
            }
        },

        /**
         * @param {jQuery} $collapsible
         */
        open: function($collapsible) {
            this.closeAll($collapsible.parent());
            $collapsible.addClass(this.options.activeClass);
            this.getContent($collapsible).slideDown();
            this.triggerEvent(eventNames.collapsibleOpened, $collapsible);
        },

        /**
         * @param {jQuery} $collapsible
         */
        close: function($collapsible) {
            $collapsible.removeClass(this.options.activeClass);
            this.getContent($collapsible).slideUp();
            this.triggerEvent(eventNames.collapsibleClosed, $collapsible);
        },

        /**
         * @param {jQuery} $collapsibleList
         */
        closeAll: function($collapsibleList) {
            var that = this;

            $collapsibleList.children().each(
                function() {
                    var $collapsible = $(this);
                    if ($collapsible.hasClass(that.options.activeClass)) {
                        that.close($collapsible);
                    }
                }
            );
        },

        /**
         * @param {string} eventName
         * @param {jQuery} $collapsible
         */
        triggerEvent: function(eventName, $collapsible) {
            this.$moduleElement.trigger(
                eventName,
                {
                    collapsible: $collapsible
                }
            );
        },

        /**
         * @param {jQuery} $collapsible
         * @returns {jQuery}
         */
        getContent: function($collapsible) {
            return $collapsible.children('[data-rd-collapsible=content]');
        }
    };

    /**
     * needs to have a DOM structure like this:
     *
     * <div data-rd-tabs="wrapper">
     *
     *      <-- tab headers, can be nested in DOM //-->
     *      <ul>
     *          <li class="rd-active" data-rd-tabs="head" data-rd-target="123">Tab1</li>
     *          <li data-rd-tabs="head" data-rd-target="789">Tab2</li>
     *      </ul>
     *
     *      <-- tabs-contents need to be a direct child of <div data-rd-tabs="wrapper"> //-->
     *      <div data-rd-tabs="contents">
     *          <div data-rd-identifier="123" style="display: block">content 123</div>
     *          <div data-rd-identifier="789">content 789</div>
     *      </div>
     * </div>
     *
     * @param {jQuery} $moduleElement
     * @param {Object} options
     * @constructor
     */
    function TabsManager($moduleElement, options) {
        this.$moduleElement = $moduleElement;
        this.options = $.extend(true, {}, this.defaultOptions, options);
        this.setEvents();
    }

    TabsManager.prototype = {
        defaultOptions: {
            activeClass: 'rd-active'
        },

        setEvents: function() {
            var that = this;

            this.$moduleElement.on('click', '[data-rd-tabs=head]', function (e) {
                e.preventDefault();
                that.onClick($(this));
            });
        },

        /**
         * @param {jQuery} $tabHead
         */
        onClick: function($tabHead) {
            if ($tabHead.hasClass(this.options.activeClass)) {
                this.triggerEvent($tabHead, eventNames.tabAlreadyShown);
                return;
            }

            this.show($tabHead);
            this.triggerEvent($tabHead, eventNames.tabShown);
        },

        /**
         * @param {jQuery} $tabHead
         */
        show: function($tabHead) {
            this.highlightHead($tabHead);
            this.showContent($tabHead);
        },

        /**
         * @param {jQuery} $tabHead
         */
        highlightHead: function($tabHead) {
            $tabHead.siblings().removeClass(this.options.activeClass);
            $tabHead.addClass(this.options.activeClass);
        },

        /**
         * @param {jQuery} $tabHead
         */
        clear: function($tabHead) {
            $tabHead.parent().children().removeClass(this.options.activeClass);
        },

        /**
         * @param {jQuery} $tabHead
         */
        showContent: function($tabHead) {
            var $contentsWrapper = this.getContentsWrapper($tabHead);

            $contentsWrapper.children().hide();
            $contentsWrapper.children('[data-rd-identifier="' + $tabHead.data('rd-target') + '"]').show();
        },

        /**
         * @param {jQuery} $tabHead
         * @param {string} eventName
         */
        triggerEvent: function($tabHead, eventName) {
            this.$moduleElement.trigger(eventName, { tabHead: $tabHead });
        },

        /**
         * @param {jQuery} $tabHead
         * @returns {jQuery}
         */
        getContentsWrapper: function($tabHead) {
            return $tabHead.closest('[data-rd-tabs=wrapper]')
                .children('[data-rd-tabs=contents]');
        }
    };

    /**
     * needs to have a DOM structure like this:
     *
     * <div data-rd-multi-selection="wrapper">
     *  <p data-rd-multi-selection="items-group" data-rd-selected="json">
     *       <a href="#" data-rd-multi-selection="item" data-rd-value="json"
     *           class="rd-active">application/json</a>
     *       <a href="#" data-rd-multi-selection="item" data-rd-value="xml">application/xml</a>
     *  </p>
     *  ... add more "items-group" as siblings if needed ....
     *  <div class="rd-code-example" data-rd-multi-selection="contents">
     *      ... build identifier with '__' as separator for each data-rd-value like 'multi-selection__json__example1'
     *      <div class="show" data-rd-identifier="multi-selection__json"></div>
     *      <div class="show" data-rd-identifier="multi-selection__xml"></div>
     *   </div>
     * </div>
     *
     * @param {jQuery} $moduleElement
     * @param {Object} options
     * @constructor
     */
    function MultiSelectionManager($moduleElement, options) {
        this.$moduleElement = $moduleElement;
        this.options = $.extend(true, {}, this.defaultOptions, options);
        this.setEvents();
    }

    MultiSelectionManager.prototype = {
        defaultOptions: {
            activeClass: 'rd-active'
        },

        setEvents: function() {
            var that = this;

            this.$moduleElement.on('click', '[data-rd-multi-selection=item]', function (e) {
                e.preventDefault();
                that.onClick($(this));
            });
        },

        /**
         * @param {jQuery} $item
         */
        onClick: function($item) {
            var $wrapper = $item.closest('[data-rd-multi-selection=wrapper]');

            this.keepSelectedItemState($item);
            this.highlightItem($item);
            this.showContent($wrapper, this.getContentIdentifier($wrapper));
        },

        /**
         * @param {jQuery} $item
         */
        keepSelectedItemState: function($item) {
            var $group = $item.closest('[data-rd-multi-selection=items-group]');
            $group.data("rd-selected",  $item.data('rd-value'));
        },

        /**
         * @param {jQuery} $item
         */
        highlightItem: function($item) {
            $item.siblings().removeClass(this.options.activeClass);
            $item.addClass(this.options.activeClass);
        },

        /**
         * @param {jQuery} $wrapper
         * @returns {string}
         */
        getContentIdentifier: function($wrapper) {
            var identifier = 'multi-selection';
            $wrapper.children('[data-rd-multi-selection=items-group]').each(
                function() {
                    identifier += '__' + $(this).data('rd-selected');
                }
            );
            return identifier;
        },

        /**
         * @param {jQuery} $wrapper
         * @param {string} identifier
         */
        showContent: function($wrapper, identifier) {
            var $contentsWrapper = $wrapper.children('[data-rd-multi-selection=contents]');
            $contentsWrapper.children().hide();
            $contentsWrapper.children('[data-rd-identifier="' + identifier + '"]').show();
        }
    };

    function ResourcesManager($moduleElement, collapsibleManager, tabsManager, options) {
        this.$moduleElement = $moduleElement;
        this.collapsibleManager = collapsibleManager;
        this.tabsManager = tabsManager;
        this.options = $.extend(true, {}, this.defaultOptions, options);

        this.$resources = $(this.options.resourcesSelector);

        this.setEventListener();
        this.setEvents();
    }

    ResourcesManager.prototype = {
        defaultOptions: {
            activeClass: 'rd-active',
            collapsibleSelector: '.rd-collapsible',
            methodsTabHeadersSelector: '> .rd-collapsible-head .rd-nested-tabs-item',
            resourcesSelector: '#resources',
            activeCollapsibleSelector: '.rd-collapsible.rd-active',
            collapsibleHeadSelector: '.rd-collapsible-head'
        },

        setEvents: function() {
            var that = this;

            this.$moduleElement.on('click', '[data-rd-resource=toggle-link]', function (e) {
                e.preventDefault();
                that.onResourceLinkClick($(this));
            });
        },

        setEventListener: function() {
            var that = this;

            this.$moduleElement.on(eventNames.tabShown, function(e, payload) {
                var $tabHead = payload.tabHead;
                if ('resource-tab' === $tabHead.data('rd-type')) {
                    that.handleTabShown($tabHead);
                    that.handleCollapsibleHeaderStyle();
                }

            });

            this.$moduleElement.on(eventNames.tabAlreadyShown, function(e, payload) {
                var $tabHead = payload.tabHead;
                if ('resource-tab' === $tabHead.data('rd-type')) {
                    that.handleTabAlreadyShown($tabHead);
                }

            });

            this.$moduleElement.on(eventNames.collapsibleOpened, function(e, payload) {
                that.handleCollapsibleHeaderStyle();
            });

            this.$moduleElement.on(eventNames.collapsibleClosed, function(e, payload) {
                var $collapsible = payload.collapsible;
                that.clearTabHeadersByCollapsible($collapsible);
                that.handleCollapsibleHeaderStyle();
            });
        },

        /**
         * @param {jQuery} $tabHead
         */
        handleTabShown: function($tabHead) {
            var $contentsWrapper = this.tabsManager.getContentsWrapper($tabHead);
            if (!$contentsWrapper.is(":visible")) {
                var $collapsible = this.getCollapsibleByChildElement($tabHead);
                this.collapsibleManager.open($collapsible);
            }
        },

        /**
         * @param {jQuery} $tabHead
         */
        handleTabAlreadyShown: function($tabHead) {
            var $collapsible = this.getCollapsibleByChildElement($tabHead);
            this.collapsibleManager.close($collapsible);

            this.tabsManager.clear($tabHead);
        },

        /**
         * @param {jQuery} $link
         */
        onResourceLinkClick: function($link) {
            var $collapsible = this.getCollapsibleByChildElement($link),
                $contentsWrapper = $collapsible.children('[data-rd-collapsible=content]'),
                $children = $contentsWrapper.children('[data-rd-resource=children]');

            if (0 == $children.length) {
                return;
            }

            var showChildren = function() {
                $contentsWrapper.children().hide();
                $children.show();
            };

            if ($collapsible.hasClass(this.options.activeClass)) {
                if ($children.is(":visible")) {
                    this.collapsibleManager.close($collapsible);
                } else {
                    this.clearTabHeadersByCollapsible($collapsible);
                    showChildren();
                }
            } else {
                showChildren();
                this.collapsibleManager.open($collapsible);
            }
        },

        /**
         * @param {jQuery} $child
         * @returns {jQuery}
         */
        getCollapsibleByChildElement: function($child) {
            return $child.closest(this.options.collapsibleSelector);
        },

        /**
         * @param {jQuery} $collapsible
         */
        clearTabHeadersByCollapsible: function($collapsible) {
            var $httpMethodsTabs = $collapsible.find(this.options.methodsTabHeadersSelector);
            if (0 < $httpMethodsTabs.length) {
                this.tabsManager.clear($httpMethodsTabs.eq(0));
            }
        },

        handleCollapsibleHeaderStyle: function() {
            var $collapsible = this.$resources.children(this.options.activeCollapsibleSelector);

            if (0 == $collapsible.length) {
                return;
            }

            var $collapsibleHeader = $collapsible.children(this.options.collapsibleHeadSelector);

            $collapsibleHeader.removeClass('rd-child-opened');

            var $childrenVisible = $collapsible.children('[data-rd-collapsible=content]')
                .children('[data-rd-resource=children]').filter(":visible");
            if ($childrenVisible.length) {
                var $childCollapsible = $childrenVisible.children().eq(0)
                    .children(this.options.activeCollapsibleSelector);
                if ($childCollapsible.length) {
                    $collapsibleHeader.addClass('rd-child-opened');
                }
            }
        }
    };

    function init($rootElement) {
        var tabsManager = new TabsManager($rootElement);
        var collapsibleManager = new CollapsibleManager($rootElement);
        new MultiSelectionManager($rootElement);
        new ResourcesManager($rootElement, collapsibleManager, tabsManager);
    }

    return {
        init: init
    };
})();

function syntaxHighlight(json) {
    json = json.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;');
    return json.replace(/("(\\u[a-zA-Z0-9]{4}|\\[^u]|[^\\"])*"(\s*:)?|\b(true|false|null)\b|-?\d+(?:\.\d*)?(?:[eE][+\-]?\d+)?)/g, function (match) {
        var cls = 'number';
        if (/^"/.test(match)) {
            if (/:$/.test(match)) {
                cls = 'key';
            } else {
                cls = 'string';
            }
        } else if (/true|false/.test(match)) {
            cls = 'boolean';
        } else if (/null/.test(match)) {
            cls = 'null';
        }
        return '<span class="' + cls + '">' + match + '</span>';
    });
}

$(document).ready(function() {
    $('.rd-json-render').each(function() {
        $(this).html(
            syntaxHighlight($(this).html())
        );
    });
});