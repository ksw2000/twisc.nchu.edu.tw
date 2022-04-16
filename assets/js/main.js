// util function
function stripHtml(html) {
    var tmp = document.createElement("div");
    tmp.innerHTML = html;
    return tmp.textContent || tmp.innerText || "";
}

function keyEnter(e, callback) {
    let keycode = (window.event) ? window.event.keyCode : e.which;
    if (keycode === 13 && callback) callback();
}

/**
 * @param {string} lang enum{'zh', 'en'}
 * @return void
 */
function switchLang(lang){
    document.cookie = "lang=" + lang;
    location.reload();
}

function renderAttachment(list) {
    ret = '';
    list.forEach(e => {
        ret += `<li><a href="${e.path}">${e.client_name}</a></li>`;
    });
    console.log(ret);
    return ret;
}

async function post(url, data) {
    return fetch(url, {
        body: data ? Object.keys(data).map(key => encodeURIComponent(key) + '=' + encodeURIComponent(data[key])).join('&') : "",
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
        },
        method: 'POST'
    }).then(res => {
        try {
            return res.json();
        } catch (e) {
            console.log(res);
        }
    });
}

function appendMoreInfo(obj) {
    $(obj).next().slideToggle();
}

const articleTypeMap = {
    'normal': '一般消息',
    'activity': '演講 & 活動',
    'course': '課程 & 招生',
    'scholarships': '獎學金',
    'recruit': '徵才資訊'
}

function notice(msg, t) {
    const dom = document.getElementById('notice');
    const timeout = t || 10000; // unit: ms
    dom.innerHTML = msg;
    dom.classList.add('show');
    setTimeout(function () {
        dom.classList.remove('show');
    }, timeout);
}

// -------------------------------- layout page --------------------------------
function openSideBar() {
    const menuBtn = document.querySelector('#button-less-900px a i');
    const sideBar = document.getElementById('button-over-900px');
    if (sideBar.classList.contains('show')) {
        sideBar.classList.remove('show');
        menuBtn.innerHTML = 'menu';
    } else {
        sideBar.classList.add('show');
        menuBtn.innerHTML = 'close';
    }
}

/**
 * @param {string} scope enum{'public', 'public-with-type'}
 * @param {string} type enum{'normal', 'activity', 'course', 'scholarships', 'recruit'}
 * @param {number|null} from? 
 * @param {number|null} to?
 * @return void 
 */
function loadNews(scope, type, from, to) {
    return new Promise((resolve, reject) => {
        $.ajax({
            url: '/api/get_news',
            data: {
                'scope': scope,
                'type': type || 'normal',
                'from': from || '',
                'to': to || ''
            },
            type: 'GET',
            success: function (data) {
                resolve(data);
            },
            error: function (err) {
                reject(err);
            },
            dataType: 'json'
        });
    });
}

// --------------------------------- login page --------------------------------
function login() {
    const id = document.querySelector("#login #id").value;
    const pwd = document.querySelector("#login #pwd").value;
    if (id === "" || pwd === "") {
        notice("帳號或密碼不可為空！");
        return;
    }
    post('/api/login', {
        id: id,
        pwd: pwd
    }).then(data => {
        if (data.err) {
            notice(data.err);
        } else {
            window.location.href = "/manage";
        }
    });
}

// ---------------------------------- reg page ---------------------------------
function register() {
    $.post('/api/reg', {
        id: $("#reg #id").val(),
        pwd: $("#reg #pwd").val(),
        re_pwd: $("#reg #re-pwd").val(),
        name: $("#reg #name").val()
    }, (data) => {
        if (data.err) {
            notice(data.err);
        } else {
            window.location.href = "/manage/reg-done";
        }
    }, 'json');
}

// ----------------------------------- routing ---------------------------------
// for rendering nav bar
$(()=>{
    try {
        let url = window.location.pathname.split('/')[1];
        let domList = $('#nav-routing-wrapper li');
        for(let i = 0; i < domList.length; i++){
            let dom = domList[i];
            if (dom.id === `nav-routing-${url}`) {
                $(dom).addClass('nav-routing-current');
            } else {
                $(dom).removeClass('nav-routing-current');
            }
        }
    } catch (e) { console.log(e) }
});

// ---------------------------------- modify pwd page ---------------------------------
function modifyPwd() {
    $.post('/api/pwd', {
        pwd: $("#modify-pwd #pwd").val(),
        re_pwd: $("#modify-pwd #re-pwd").val()
    }, (data) => {
        if (data.err) {
            notice(data.err);
        } else {
            window.location.href = "/manage/pwd-done";
        }
    }, 'json');
}