import {h, app} from '/static/3rd/ha.js'

// effects
const _fetchFx = (dispatch, {path, action, options, args}) => {
    fetch(path, options)
        .then(response => {
            if (response.ok) return response.json()

            return response.text().then(body => {
                throw new Error(response.statusText + ': ' + body)
            })
        })
        .then(data => dispatch(action, {data, args}))
        .catch(err => dispatch(action, {err, args}))
}

// actions
const Unlock = state => ({...state, lock: state.lock - 1})
const Lock = state => ({...state, lock: state.lock + 1})
const UpdateNewUser = (state, e) => ({...state, newUser: e.target.value})
const UpdateNewEmail = (state, e) => ({
    ...state,
    newEmail: {...state.newEmail, address: e.target.value},
})

const handleFetchEmails = (state, {data = [], err, args}) => {
    if (err) {
        alert('could not fetch e-mails')
        console.error(err)
        return Unlock
    }

    return Unlock({...state, emails: [...state.emails, ...data], newEmail: {user: args}})
}
const FetchEmails = (state, id) => [
    Lock(state),
    [_fetchFx, {action: handleFetchEmails, path: '/rest/emails/' + id, args: id}],
]

const handleFetchUsers = (state, {data=[], err}) => {
    if (err) {
        alert('could fetch users')
        console.error(err)
        return Unlock
    }

    return Unlock({...state, users: data, loading: null})
}
const FetchUsers = (state) => [
    Lock({...state, loading: 'users'}),
    [_fetchFx, {action: handleFetchUsers, path: '/rest/users'}],
]

const handleAddUser = (state, {data, err}) => {
    if (err) {
        alert('could not add user')
        console.error(err)
        return Unlock
    }

    return [
        Unlock({...state, newUser: ''}),
        [(d, {action}) => d(action), {action: FetchUsers}],
    ]
}
const AddUser = state => [
    Lock(state),
    [
        _fetchFx,
        {
            action: handleAddUser,
            path: '/rest/users',
            options: {
                method: 'POST',
                body: JSON.stringify({name: state.newUser}),
            },
        },
    ],
]

const handleRemoveUser = (state, {data, err, args}) => {
    if (err) {
        alert('could not remove user')
        console.error(err)
    }

    return Unlock({...state, users: state.users.filter(u => u.id != args)})
}
const RemoveUser = (state, userID) => [
    Lock(state),
    [
        _fetchFx,
        {
            args: userID,
            action: handleRemoveUser,
            path: '/rest/users',
            options: {
                method: 'DELETE',
                body: JSON.stringify({userID}),
            },
        },
    ],
]

const handleAddEmail = (state, {data, err}) => {
    if (err) {
        alert('could not add e-mail address')
        console.error(err)
        return Unlock
    }

    return [
        Unlock,
        [(d, {action}) => d(action), {action: [FetchEmails, state.newEmail.user]}],
    ]
}
const AddEmail = state => [
    Lock(state),
    [
        _fetchFx,
        {
            action: handleAddEmail,
            path: '/rest/emails?debug',
            options: {
                method: 'POST',
                body: JSON.stringify(state.newEmail),
            },
        },
    ],
]

const User = (state, user) => h(
    'li',
    {key: user.id},
    state.newEmail && state.newEmail.user == user.id && h(
        'div', {className:'card'},
        [
            h('header', {className: 'card-header'}, [
                h('p', {className: 'card-header-title'}, user.name),
                h('a', {className:'card-header-icon', disabled: state.lock, onclick: [RemoveUser, user.id]},
                    h('span', {className: 'delete'}),
                ),
            ]),
            h('div', {className: 'card-content'}, [
                h('div', {className: 'content'}, [
                    h(
                        'div',
                        {className: 'field has-addons'},
                        h('div', {className: 'control'}, h('input', {
                            className:'input',
                            type: 'text',
                            disabled: state.lock,
                            value: state.newEmail.address,
                            oninput: UpdateNewEmail,
                            placeHolder: 'Email',
                        })),
                        h('div', {className: 'control'}, h('button', {
                            className:'button is-primary', disabled: state.lock || !state.newEmail.address, onclick: AddEmail,
                        }, '+'))
                    ),
                    state.emails
                        .filter(e => e.user_id == user.id)
                        .map(e => h('div', null, e.address, h('a', {className: 'delete is-small'}))),
                ]),
            ]),
        ]
    ) || h(
        'strong',
        {
            'data-id': user.id,
            'onclick': [FetchEmails, user.id],
        },
        user.name
    ),
)

app({
    init: () => ({
        lock: 0,
        loading: null,
        newUser: '',
        newEmail: {},
        emails: [],
        users: [],
    }),
    view: state =>
        h('div', {className: 'section'}, [
            h('div', {className: 'container'}, [
                h('h1', {className: 'title'}, 'Rest'),
                h('a', {href: '/graphql/play'}, 'graphql'),
                h(
                    'div',
                    {className: 'field has-addons'},
                    h(
                        'div',
                        {className:'control'},
                        h(
                            'button',
                            {
                                class:{'button': true, 'is-loading': state.loading == 'users'},
                                onclick: FetchUsers, disabled: state.lock,
                            },
                            'fetch users'
                        )
                    ),
                    h('div', {className: 'control'},
                        h('input', {
                            className: 'input',
                            type: 'text',
                            placeHolder: 'Name',
                            disabled: state.lock,
                            oninput: UpdateNewUser,
                            value: state.newUser,
                        })),
                    h('div', {className: 'control'},
                        h('button', {className:'button is-primary',
                            disabled: state.lock || !state.newUser, onclick: AddUser}, 'add user')
                    )
                ),
                h(
                    'ul',
                    null,
                    state.users.map(u => User(state, u))
                ),
            ]),
        ]),
    node: document.getElementById('app'),
})
