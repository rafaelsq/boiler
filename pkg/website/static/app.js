import {h, app} from '/static/3rd/ha.js'

const UpdateUsers = (state, users) => ({...state, users})
const UpdateEmails = (state, emails) => ({...state, emails})

const _fetchFx = ({path, action}, dispatch) =>
    fetch(path)
        .then(data => {
            return data.json()
        })
        .then(data => {
            dispatch(action, data)
        })

const FetchEmails = (state, id) => [
    state,
    [_fetchFx, {action: UpdateEmails, path: '/rest/emails/' + id}],
]

const FetchUsers = state => [
    state,
    [_fetchFx, {action: UpdateUsers, path: '/rest/users'}],
]

app({
    init: () => ({
        emails: [],
        users: [],
    }),
    view: state =>
        h('div', {}, [
            h('h1', {}, 'Rest'),
            h('a', {href: '/graphql/play'}, 'graphql'),
            h('button', {onclick: FetchUsers}, 'fetch users'),
            state.users.map(u =>
                h(
                    'li',
                    {
                        'data-id': u.id,
                        'onclick': [FetchEmails, u.id],
                    },
                    u.name
                )
            ),
            state.emails.map(e => h('ol', {}, e.address)),
        ]),
    node: document.getElementById('app'),
})
