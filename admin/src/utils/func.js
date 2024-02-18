import dayjs from 'dayjs'
export function timeDifference(toDate) {
    // 当前日期和时间
    var now = new Date();

    // 将输入的日期字符串转换为 Date 对象
    var targetDate = new Date(toDate.replace(/-/g, '/')); // 替换破折号以兼容不同浏览器

    // 计算时间差（毫秒）
    var diff = targetDate - now;

    // 如果时间已经过去
    if (diff < 0) {
        return "已过期";
    }

    // 转换为天、小时、分钟
    var days = Math.floor(diff / (1000 * 60 * 60 * 24));
    var hours = Math.floor((diff % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
    var minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60));

    return `${days ?  days + '天' : ''} ${hours ? hours + '小时' : ''} ${minutes} 分`;
}

export function dateFormat(toDate) {
    return dayjs(toDate).format('YYYY-MM-DD HH:mm:ss')
}
